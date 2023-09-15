package assured

import (
	"bytes"
	"context"
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/go-kit/kit/endpoint"
)

// AssuredEndpoints
type AssuredEndpoints struct {
	httpClient     *http.Client
	assuredCalls   *CallStore
	madeCalls      *CallStore
	callbackCalls  *CallStore
	trackMadeCalls bool
}

// NewAssuredEndpoints creates a new instance of assured endpoints
func NewAssuredEndpoints(options Options) *AssuredEndpoints {
	return &AssuredEndpoints{
		assuredCalls:   NewCallStore(),
		madeCalls:      NewCallStore(),
		callbackCalls:  NewCallStore(),
		httpClient:     options.httpClient,
		trackMadeCalls: options.trackMadeCalls,
	}
}

// WrappedEndpoint is used to validate that the incoming request is an assured call
func (a *AssuredEndpoints) WrappedEndpoint(handler func(context.Context, *Call) (interface{}, error)) endpoint.Endpoint {
	return func(ctx context.Context, i interface{}) (response interface{}, err error) {
		a, ok := i.(*Call)
		if !ok {
			return nil, errors.New("unable to convert request to assured Call")
		}

		return handler(ctx, a)
	}
}

// GivenEndpoint is used to stub out a call for a given path
func (a *AssuredEndpoints) GivenEndpoint(ctx context.Context, call *Call) (interface{}, error) {
	a.assuredCalls.Add(call)
	slog.With("path", call.ID()).Info("assured call set")

	return call, nil
}

// GivenCallbackEndpoint is used to stub out callbacks for a callback key
func (a *AssuredEndpoints) GivenCallbackEndpoint(ctx context.Context, call *Call) (interface{}, error) {
	a.callbackCalls.AddAt(call.Headers[AssuredCallbackKey], call)
	slog.With("key", call.Headers[AssuredCallbackKey], "target", call.Headers[AssuredCallbackTarget]).Info("assured callback set")

	return call, nil
}

// WhenEndpoint is used to test the assured calls
func (a *AssuredEndpoints) WhenEndpoint(ctx context.Context, call *Call) (interface{}, error) {
	calls := a.assuredCalls.Get(call.ID())
	if len(calls) == 0 {
		slog.With("path", call.ID()).Info("assured call not found")
		return nil, errors.New("No assured calls")
	}

	if a.trackMadeCalls {
		a.madeCalls.Add(call)
	}
	assured := calls[0]
	a.assuredCalls.Rotate(assured)

	// Trigger callbacks, if applicable
	for _, callback := range a.callbackCalls.Get(assured.Headers[AssuredCallbackKey]) {
		go a.sendCallback(callback.Headers[AssuredCallbackTarget], callback)
	}

	// Delay response
	if delay, err := strconv.ParseInt(assured.Headers[AssuredDelay], 10, 64); err == nil {
		time.Sleep(time.Duration(delay) * time.Second)
	}

	slog.With("path", call.ID()).Info("assured call responded")
	return assured, nil
}

// VerifyEndpoint is used to verify a particular call
func (a *AssuredEndpoints) VerifyEndpoint(ctx context.Context, call *Call) (interface{}, error) {
	if a.trackMadeCalls {
		return a.madeCalls.Get(call.ID()), nil
	}
	return nil, errors.New("Tracking made calls is disabled")
}

// ClearEndpoint is used to clear a specific assured call
func (a *AssuredEndpoints) ClearEndpoint(ctx context.Context, call *Call) (interface{}, error) {
	a.assuredCalls.Clear(call.ID())
	a.madeCalls.Clear(call.ID())
	slog.With("path", call.ID()).Info("cleared calls for path")
	if call.Headers[AssuredCallbackKey] != "" {
		a.callbackCalls.Clear(call.Headers[AssuredCallbackKey])
		slog.With("key", call.Headers[AssuredCallbackKey]).Info("cleared calls for key")
	}

	return nil, nil
}

// ClearAllEndpoint is used to clear all assured calls
func (a *AssuredEndpoints) ClearAllEndpoint(ctx context.Context, i interface{}) (interface{}, error) {
	a.assuredCalls.ClearAll()
	a.madeCalls.ClearAll()
	a.callbackCalls.ClearAll()
	slog.Info("cleared all calls")

	return nil, nil
}

// sendCallback sends a given callback to its target
func (a *AssuredEndpoints) sendCallback(target string, call *Call) {
	var delay int64
	if delayOverride, err := strconv.ParseInt(call.Headers[AssuredCallbackDelay], 10, 64); err == nil {
		delay = delayOverride
	}
	req, err := http.NewRequest(call.Method, target, bytes.NewBuffer(call.Response))
	if err != nil {
		slog.With("target", target, "error", err).Info("failed to build callback request")
		return
	}
	for key, value := range call.Headers {
		req.Header.Set(key, value)
	}
	// Delay callback, if applicable
	time.Sleep(time.Duration(delay) * time.Second)
	resp, err := a.httpClient.Do(req)
	if err != nil {
		slog.With("target", target, "error", err).Info("failed to reach callback target")
		return
	}
	slog.With("target", target, "status_code", resp.StatusCode).Info("sent callback to target")
}

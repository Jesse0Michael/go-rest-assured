package assured

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/go-kit/kit/endpoint"
	kitlog "github.com/go-kit/kit/log"
)

// AssuredEndpoints
type AssuredEndpoints struct {
	logger         kitlog.Logger
	httpClient     http.Client
	assuredCalls   *CallStore
	madeCalls      *CallStore
	trackMadeCalls bool
}

// Settings
type Settings struct {
	Logger         kitlog.Logger
	HTTPClient     http.Client
	Port           int
	TrackMadeCalls bool
}

// NewAssuredEndpoints creates a new instance of assured endpoints
func NewAssuredEndpoints(settings Settings) *AssuredEndpoints {
	return &AssuredEndpoints{
		assuredCalls:   NewCallStore(),
		madeCalls:      NewCallStore(),
		logger:         settings.Logger,
		httpClient:     settings.HTTPClient,
		trackMadeCalls: settings.TrackMadeCalls,
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
	// Assign Call as callback, if applicable
	if call.Headers[AssuredCallbackKey] != "" && call.Headers[AssuredCallbackTarget] != "" {
		changed := a.assuredCalls.AddCallback(call.Headers[AssuredCallbackKey], call)
		if len(changed) == 0 {
			a.logger.Log("message", "assured callback key not found", "path", call.Headers[AssuredCallbackKey])
			return nil, errors.New("No assured callback key found")
		}
		a.logger.Log("message", "assured callback set", "target", call.Headers[AssuredCallbackTarget])

		return call, nil
	}

	a.assuredCalls.Add(call)
	a.logger.Log("message", "assured call set", "path", call.ID())

	return call, nil
}

// WhenEndpoint is used to test the assured calls
func (a *AssuredEndpoints) WhenEndpoint(ctx context.Context, call *Call) (interface{}, error) {
	calls := a.assuredCalls.Get(call.ID())
	if len(calls) == 0 {
		a.logger.Log("message", "assured call not found", "path", call.ID())
		return nil, errors.New("No assured calls")
	}

	if a.trackMadeCalls {
		a.madeCalls.Add(call)
	}
	assured := calls[0]
	a.assuredCalls.Rotate(assured)

	// Trigger callbacks, if applicable
	for _, callback := range assured.Callbacks {
		go a.sendCallback(callback.Headers[AssuredCallbackTarget], callback)
	}

	a.logger.Log("message", "assured call responded", "path", call.ID())
	return assured, nil
}

// VerifyEndpoint is used to verify a particular call
func (a *AssuredEndpoints) VerifyEndpoint(ctx context.Context, call *Call) (interface{}, error) {
	if a.trackMadeCalls {
		return a.madeCalls.Get(call.ID()), nil
	}
	return nil, errors.New("Tracking made calls is disabled")
}

//ClearEndpoint is used to clear a specific assured call
func (a *AssuredEndpoints) ClearEndpoint(ctx context.Context, call *Call) (interface{}, error) {
	a.assuredCalls.Clear(call.ID())
	a.madeCalls.Clear(call.ID())
	a.logger.Log("message", "cleared calls for path", "path", call.ID())

	return nil, nil
}

//ClearAllEndpoint is used to clear all assured calls
func (a *AssuredEndpoints) ClearAllEndpoint(ctx context.Context, i interface{}) (interface{}, error) {
	a.assuredCalls.ClearAll()
	a.madeCalls.ClearAll()
	a.logger.Log("message", "cleared all calls")

	return nil, nil
}

//sendCallback sends a given callback to its target
func (a *AssuredEndpoints) sendCallback(target string, call Call) {
	var delay int64
	if delayOverride, err := strconv.ParseInt(call.Headers[AssuredCallbackDelay], 10, 64); err == nil {
		delay = delayOverride
	}
	req, err := http.NewRequest(call.Method, target, bytes.NewBuffer(call.Response))
	if err != nil {
		a.logger.Log("message", "failed to build callback request", "target", target, "error", err.Error())
	}
	for key, value := range call.Headers {
		req.Header.Set(key, value)
	}

	// Delay callback, if applicable
	time.Sleep(time.Duration(delay) * time.Second)
	resp, err := a.httpClient.Do(req)
	if err != nil {
		a.logger.Log("message", "failed to reach callback target", "target", target, "error", err.Error())
	}
	a.logger.Log("message", "sent callback to target", "target", target, "status_code", resp.StatusCode)
}

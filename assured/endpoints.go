package assured

import (
	"context"
	"errors"

	"github.com/go-kit/kit/endpoint"
	kitlog "github.com/go-kit/kit/log"
)

// AssuredEndpoints
type AssuredEndpoints struct {
	logger         kitlog.Logger
	assuredCalls   *CallStore
	madeCalls      *CallStore
	trackMadeCalls bool
}

// Settings
type Settings struct {
	Logger         kitlog.Logger
	Port           int
	TrackMadeCalls bool
}

// NewAssuredEndpoints creates a new instance of assured endpoints
func NewAssuredEndpoints(settings Settings) *AssuredEndpoints {
	return &AssuredEndpoints{
		assuredCalls:   NewCallStore(),
		madeCalls:      NewCallStore(),
		logger:         settings.Logger,
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

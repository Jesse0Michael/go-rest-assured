package endpoints

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-kit/kit/endpoint"
	kitlog "github.com/go-kit/kit/log"
	"github.com/jesse0michael/go-rest-assured/assured"
)

// AssuredEndpoints
type AssuredEndpoints struct {
	logger       kitlog.Logger
	assuredCalls map[string][]*assured.Call
	madeCalls    map[string][]*assured.Call
}

// NewAssuredEndpoints creates a new instance of assured endpoints
func NewAssuredEndpoints(l kitlog.Logger) *AssuredEndpoints {
	return &AssuredEndpoints{
		logger:       l,
		assuredCalls: map[string][]*assured.Call{},
		madeCalls:    map[string][]*assured.Call{},
	}
}

// WrappedEndpoint is used to validate that the incoming request is an assured call
func (a *AssuredEndpoints) WrappedEndpoint(handler func(context.Context, *assured.Call) (interface{}, error)) endpoint.Endpoint {
	return func(ctx context.Context, i interface{}) (response interface{}, err error) {
		a, ok := i.(*assured.Call)
		if !ok {
			return nil, errors.New("unable to convert request to assured Call")
		}
		return handler(ctx, a)
	}
}

// GivenEndpoint is used to stub out a call for a given path
func (a *AssuredEndpoints) GivenEndpoint(ctx context.Context, call *assured.Call) (interface{}, error) {
	if a.assuredCalls[call.ID()] == nil {
		a.assuredCalls[call.ID()] = []*assured.Call{call}
	} else {
		a.assuredCalls[call.ID()] = append(a.assuredCalls[call.ID()], call)
	}
	a.logger.Log("message", "assured call set", "path", call.ID())

	return call, nil
}

// WhenEndpoint is used to test the assured calls
func (a *AssuredEndpoints) WhenEndpoint(ctx context.Context, call *assured.Call) (interface{}, error) {
	if a.assuredCalls[call.ID()] == nil || len(a.assuredCalls[call.ID()]) == 0 {
		a.logger.Log("message", "assured call not found", "path", call.ID())
		return nil, errors.New("No assured calls")
	}

	a.madeCalls[call.ID()] = append(a.madeCalls[call.ID()], call)

	assured := a.assuredCalls[call.ID()][0]

	a.assuredCalls[call.ID()] = append(a.assuredCalls[call.ID()][1:], assured)

	return assured, nil
}

// ThenEndpoint is used to verify a particular call
func (a *AssuredEndpoints) ThenEndpoint(ctx context.Context, call *assured.Call) (interface{}, error) {
	return a.madeCalls[call.ID()], nil
}

//ClearEndpoint is used to clear a specific assured call
func (a *AssuredEndpoints) ClearEndpoint(ctx context.Context, call *assured.Call) (interface{}, error) {
	delete(a.assuredCalls, call.ID())
	delete(a.madeCalls, call.ID())
	a.logger.Log("message", "cleared calls for path", "path", call.ID())

	return nil, nil
}

//ClearAllEndpoint is used to clear all assured calls
func (a *AssuredEndpoints) ClearAllEndpoint(ctx context.Context, i interface{}) (interface{}, error) {
	a.assuredCalls = map[string][]*assured.Call{}
	a.madeCalls = map[string][]*assured.Call{}
	a.logger.Log("message", "cleared all calls")

	return nil, nil
}

package endpoints

import (
	"context"
	"errors"

	kitlog "github.com/go-kit/kit/log"
)

type AssuredCall struct {
	Path       string
	Response   interface{}
	StatusCode int
}

// AssuredEndpoints
type AssuredEndpoints struct {
	logger       kitlog.Logger
	assuredCalls map[string][]*AssuredCall
	madeCalls    map[string][]*AssuredCall
}

// NewAssuredEndpoints creates a new instance of assured endpoints
func NewAssuredEndpoints(l kitlog.Logger) *AssuredEndpoints {
	return &AssuredEndpoints{
		logger:       l,
		assuredCalls: map[string][]*AssuredCall{},
		madeCalls:    map[string][]*AssuredCall{},
	}
}

// GivenEndpoint is used to stub out a call for a given path
func (a *AssuredEndpoints) GivenEndpoint(ctx context.Context, i interface{}) (interface{}, error) {
	call, ok := i.(*AssuredCall)
	if !ok {
		return nil, errors.New("unable to convert request to AssuredCall")
	}
	a.assuredCalls[call.Path] = append(a.assuredCalls[call.Path], call)
	a.logger.Log("message", "assured call set", "path", call.Path)

	return call, nil
}

// WhenEndpoint is used to test the assured calls
func (a *AssuredEndpoints) WhenEndpoint(ctx context.Context, i interface{}) (interface{}, error) {
	call, ok := i.(*AssuredCall)
	if !ok {
		return nil, errors.New("unable to convert request to AssuredCall")
	}
	if a.assuredCalls[call.Path] == nil || len(a.assuredCalls[call.Path]) == 0 {
		a.logger.Log("message", "assured call not found", "path", call.Path)
		return nil, errors.New("No assured calls")
	}

	a.madeCalls[call.Path] = append(a.madeCalls[call.Path], call)

	assured := a.assuredCalls[call.Path][0]

	a.assuredCalls[call.Path] = append(a.assuredCalls[call.Path][1:], assured)

	return assured, nil
}

// ThenEndpoint is used to verify a particular call
func (a *AssuredEndpoints) ThenEndpoint(ctx context.Context, i interface{}) (interface{}, error) {
	call, ok := i.(*AssuredCall)
	if !ok {
		return nil, errors.New("unable to convert request to AssuredCall")
	}

	return a.madeCalls[call.Path], nil
}

//ClearEndpoint is used to clear a specific assured call
func (a *AssuredEndpoints) ClearEndpoint(ctx context.Context, i interface{}) (interface{}, error) {
	call, ok := i.(*AssuredCall)
	if !ok {
		return nil, errors.New("unable to convert request to AssuredCall")
	}
	delete(a.assuredCalls, call.Path)
	delete(a.madeCalls, call.Path)
	a.logger.Log("message", "cleared calls for path", "path", call.Path)

	return nil, nil
}

//ClearAllEndpoint is used to clear all assured calls
func (a *AssuredEndpoints) ClearAllEndpoint(ctx context.Context, i interface{}) (interface{}, error) {
	a.assuredCalls = map[string][]*AssuredCall{}
	a.madeCalls = map[string][]*AssuredCall{}
	a.logger.Log("message", "cleared all calls")

	return nil, nil
}

package endpoints

import (
	"context"
	"errors"

	kitlog "github.com/go-kit/kit/log"
)

type AssuredCall struct {
	Path       string
	Response   map[string]interface{}
	StatusCode int
}

type AssuredEndpoints struct {
	logger       kitlog.Logger
	assuredCalls map[string]*AssuredCall
	madeCalls    map[string]*AssuredCall
}

func NewAssuredEndpoints(l kitlog.Logger) *AssuredEndpoints {
	return &AssuredEndpoints{
		logger:       l,
		assuredCalls: map[string]*AssuredCall{},
		madeCalls:    map[string]*AssuredCall{},
	}
}

func (a AssuredEndpoints) GivenEndpoint(ctx context.Context, call *AssuredCall) (*AssuredCall, error) {
	a.assuredCalls[call.Path] = call
	a.logger.Log("assured call for path: %s", call.Path)

	return call, nil
}

func (a AssuredEndpoints) WhenEndpoint(ctx context.Context, call *AssuredCall) (*AssuredCall, error) {
	if a.assuredCalls[call.Path] == nil {
		a.logger.Log("assured call not found for path: %s", call.Path)
		return nil, errors.New("No assured calls")
	}

	a.madeCalls[call.Path] = call

	return a.assuredCalls[call.Path], nil
}

func (a AssuredEndpoints) ThenEndpoint(ctx context.Context, call *AssuredCall) (*AssuredCall, error) {
	if a.madeCalls[call.Path] == nil {
		a.logger.Log("made call not found for path: %s", call.Path)
		return nil, errors.New("No made calls")
	}

	return a.madeCalls[call.Path], nil
}

func (a AssuredEndpoints) ClearEndpoint(ctx context.Context, call *AssuredCall) (*AssuredCall, error) {
	a.assuredCalls[call.Path] = nil
	a.madeCalls[call.Path] = nil
	a.logger.Log("Cleared calls for path: %s", call.Path)

	return nil, nil
}

func (a AssuredEndpoints) ClearAllEndpoint(ctx context.Context, i interface{}) (interface{}, error) {
	a.assuredCalls = map[string]*AssuredCall{}
	a.madeCalls = map[string]*AssuredCall{}
	a.logger.Log("Cleared all calls")

	return nil, nil
}

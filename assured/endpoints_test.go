package assured

import (
	"context"
	"io/ioutil"
	"net/http"
	"testing"

	kitlog "github.com/go-kit/kit/log"
	"github.com/stretchr/testify/require"
)

func TestNewAssuredEndpoints(t *testing.T) {
	logger := kitlog.NewLogfmtLogger(ioutil.Discard)
	expected := &AssuredEndpoints{
		logger:       logger,
		assuredCalls: NewCallStore(),
		madeCalls:    NewCallStore(),
	}
	actual := NewAssuredEndpoints(logger)

	require.Equal(t, expected.assuredCalls, actual.assuredCalls)
	require.Equal(t, expected.madeCalls, actual.madeCalls)
}

func TestWrappedEndpointSuccess(t *testing.T) {
	endpoints := NewAssuredEndpoints(kitlog.NewLogfmtLogger(ioutil.Discard))
	testEndpoint := func(ctx context.Context, call *Call) (interface{}, error) {
		return call, nil
	}

	actual := endpoints.WrappedEndpoint(testEndpoint)
	c, err := actual(ctx, call1)

	require.NoError(t, err)
	require.Equal(t, call1, c)
}

func TestWrappedEndpointFailure(t *testing.T) {
	endpoints := NewAssuredEndpoints(kitlog.NewLogfmtLogger(ioutil.Discard))
	testEndpoint := func(ctx context.Context, call *Call) (interface{}, error) {
		return call, nil
	}

	actual := endpoints.WrappedEndpoint(testEndpoint)
	c, err := actual(ctx, false)

	require.Nil(t, c)
	require.Error(t, err)
	require.Equal(t, err.Error(), "unable to convert request to assured Call")
}

func TestGivenEndpointSuccess(t *testing.T) {
	endpoints := NewAssuredEndpoints(kitlog.NewLogfmtLogger(ioutil.Discard))

	c, err := endpoints.GivenEndpoint(ctx, call1)

	require.NoError(t, err)
	require.Equal(t, call1, c)

	c, err = endpoints.GivenEndpoint(ctx, call2)

	require.NoError(t, err)
	require.Equal(t, call2, c)

	c, err = endpoints.GivenEndpoint(ctx, call3)

	require.NoError(t, err)
	require.Equal(t, call3, c)

	require.Equal(t, fullAssuredCalls, endpoints.assuredCalls)
}

func TestWhenEndpointSuccess(t *testing.T) {
	endpoints := &AssuredEndpoints{
		assuredCalls: fullAssuredCalls,
		madeCalls:    NewCallStore(),
	}
	expected := map[string][]*Call{
		"GET:test/assured":    {call2, call1},
		"POST:teapot/assured": {call3},
	}

	c, err := endpoints.WhenEndpoint(ctx, call1)

	require.NoError(t, err)
	require.Equal(t, call1, c)
	require.Equal(t, expected, endpoints.assuredCalls.data)

	c, err = endpoints.WhenEndpoint(ctx, call2)

	require.NoError(t, err)
	require.Equal(t, call2, c)
	require.Equal(t, fullAssuredCalls, endpoints.assuredCalls)

	c, err = endpoints.WhenEndpoint(ctx, call3)

	require.NoError(t, err)
	require.Equal(t, call3, c)
	require.Equal(t, fullAssuredCalls, endpoints.assuredCalls)
	require.Equal(t, fullAssuredCalls, endpoints.madeCalls)
}

func TestWhenEndpointNotFound(t *testing.T) {
	endpoints := NewAssuredEndpoints(kitlog.NewLogfmtLogger(ioutil.Discard))

	c, err := endpoints.WhenEndpoint(ctx, call1)

	require.Nil(t, c)
	require.Error(t, err)
	require.Equal(t, "No assured calls", err.Error())
}

func TestVerifyEndpointSuccess(t *testing.T) {
	endpoints := &AssuredEndpoints{
		madeCalls: fullAssuredCalls,
	}

	c, err := endpoints.VerifyEndpoint(ctx, call1)

	require.NoError(t, err)
	require.Equal(t, []*Call{call1, call2}, c)

	c, err = endpoints.VerifyEndpoint(ctx, call3)

	require.NoError(t, err)
	require.Equal(t, []*Call{call3}, c)
}

func TestClearEndpointSuccess(t *testing.T) {
	endpoints := &AssuredEndpoints{
		logger:       kitlog.NewLogfmtLogger(ioutil.Discard),
		assuredCalls: fullAssuredCalls,
		madeCalls:    fullAssuredCalls,
	}
	expected := map[string][]*Call{
		"POST:teapot/assured": {call3},
	}

	c, err := endpoints.ClearEndpoint(ctx, call1)

	require.NoError(t, err)
	require.Nil(t, c)
	require.Equal(t, expected, endpoints.assuredCalls.data)
	require.Equal(t, expected, endpoints.madeCalls.data)

	c, err = endpoints.ClearEndpoint(ctx, call2)

	require.NoError(t, err)
	require.Nil(t, c)
	require.Equal(t, expected, endpoints.assuredCalls.data)
	require.Equal(t, expected, endpoints.madeCalls.data)

	c, err = endpoints.ClearEndpoint(ctx, call3)

	require.NoError(t, err)
	require.Nil(t, c)
	require.Equal(t, map[string][]*Call{}, endpoints.assuredCalls.data)
	require.Equal(t, map[string][]*Call{}, endpoints.madeCalls.data)
}

func TestClearAllEndpointSuccess(t *testing.T) {
	endpoints := &AssuredEndpoints{
		logger:       kitlog.NewLogfmtLogger(ioutil.Discard),
		assuredCalls: fullAssuredCalls,
		madeCalls:    fullAssuredCalls,
	}

	c, err := endpoints.ClearAllEndpoint(ctx, nil)

	require.NoError(t, err)
	require.Nil(t, c)
	require.Equal(t, map[string][]*Call{}, endpoints.assuredCalls.data)
	require.Equal(t, map[string][]*Call{}, endpoints.madeCalls.data)
}

var (
	call1 = &Call{
		Path:       "test/assured",
		Method:     "GET",
		StatusCode: http.StatusOK,
		Response:   []byte(`{"assured": true}`),
	}
	call2 = &Call{
		Path:       "test/assured",
		Method:     "GET",
		StatusCode: http.StatusConflict,
		Response:   []byte("error"),
	}
	call3 = &Call{
		Path:       "teapot/assured",
		Method:     "POST",
		StatusCode: http.StatusTeapot,
	}
	fullAssuredCalls = &CallStore{
		data: map[string][]*Call{
			"GET:test/assured":    {call1, call2},
			"POST:teapot/assured": {call3},
		},
	}
)

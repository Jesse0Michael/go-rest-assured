package assured

import (
	"context"
	"io/ioutil"
	"testing"

	kitlog "github.com/go-kit/kit/log"
	"github.com/stretchr/testify/require"
)

func TestNewAssuredEndpoints(t *testing.T) {
	expected := &AssuredEndpoints{
		logger:         testSettings.Logger,
		assuredCalls:   NewCallStore(),
		madeCalls:      NewCallStore(),
		trackMadeCalls: true,
	}
	actual := NewAssuredEndpoints(testSettings)

	require.Equal(t, expected.assuredCalls, actual.assuredCalls)
	require.Equal(t, expected.madeCalls, actual.madeCalls)
}

func TestWrappedEndpointSuccess(t *testing.T) {
	endpoints := NewAssuredEndpoints(testSettings)
	testEndpoint := func(ctx context.Context, call *Call) (interface{}, error) {
		return call, nil
	}

	actual := endpoints.WrappedEndpoint(testEndpoint)
	c, err := actual(ctx, call1)

	require.NoError(t, err)
	require.Equal(t, call1, c)
}

func TestWrappedEndpointFailure(t *testing.T) {
	endpoints := NewAssuredEndpoints(testSettings)
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
	endpoints := NewAssuredEndpoints(testSettings)

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
		logger:         testSettings.Logger,
		assuredCalls:   fullAssuredCalls,
		madeCalls:      NewCallStore(),
		trackMadeCalls: true,
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

func TestWhenEndpointSuccessTrackingDisabled(t *testing.T) {
	endpoints := &AssuredEndpoints{
		logger:         testSettings.Logger,
		assuredCalls:   fullAssuredCalls,
		madeCalls:      NewCallStore(),
		trackMadeCalls: false,
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
	require.Equal(t, NewCallStore(), endpoints.madeCalls)
}

func TestWhenEndpointNotFound(t *testing.T) {
	endpoints := NewAssuredEndpoints(testSettings)

	c, err := endpoints.WhenEndpoint(ctx, call1)

	require.Nil(t, c)
	require.Error(t, err)
	require.Equal(t, "No assured calls", err.Error())
}

func TestVerifyEndpointSuccess(t *testing.T) {
	endpoints := &AssuredEndpoints{
		madeCalls:      fullAssuredCalls,
		trackMadeCalls: true,
	}

	c, err := endpoints.VerifyEndpoint(ctx, call1)

	require.NoError(t, err)
	require.Equal(t, []*Call{call1, call2}, c)

	c, err = endpoints.VerifyEndpoint(ctx, call3)

	require.NoError(t, err)
	require.Equal(t, []*Call{call3}, c)
}

func TestVerifyEndpointTrackingDisabled(t *testing.T) {
	endpoints := &AssuredEndpoints{
		madeCalls:      fullAssuredCalls,
		trackMadeCalls: false,
	}

	c, err := endpoints.VerifyEndpoint(ctx, call1)

	require.Nil(t, c)
	require.Error(t, err)
	require.Equal(t, "Tracking made calls is disabled", err.Error())
}

func TestClearEndpointSuccess(t *testing.T) {
	endpoints := &AssuredEndpoints{
		logger:         kitlog.NewLogfmtLogger(ioutil.Discard),
		assuredCalls:   fullAssuredCalls,
		madeCalls:      fullAssuredCalls,
		trackMadeCalls: true,
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
		logger:         kitlog.NewLogfmtLogger(ioutil.Discard),
		assuredCalls:   fullAssuredCalls,
		madeCalls:      fullAssuredCalls,
		trackMadeCalls: true,
	}

	c, err := endpoints.ClearAllEndpoint(ctx, nil)

	require.NoError(t, err)
	require.Nil(t, c)
	require.Equal(t, map[string][]*Call{}, endpoints.assuredCalls.data)
	require.Equal(t, map[string][]*Call{}, endpoints.madeCalls.data)
}

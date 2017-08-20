package endpoints

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
		assuredCalls: map[string][]*AssuredCall{},
		madeCalls:    map[string][]*AssuredCall{},
	}
	actual := NewAssuredEndpoints(logger)

	require.Equal(t, expected, actual)
}

func TestGivenEndpointSuccess(t *testing.T) {
	assured := NewAssuredEndpoints(kitlog.NewLogfmtLogger(ioutil.Discard))

	c, err := assured.GivenEndpoint(ctx, call1)

	require.NoError(t, err)
	require.Equal(t, call1, c)

	c, err = assured.GivenEndpoint(ctx, call2)

	require.NoError(t, err)
	require.Equal(t, call2, c)

	c, err = assured.GivenEndpoint(ctx, call3)

	require.NoError(t, err)
	require.Equal(t, call3, c)

	require.Equal(t, fullAssuredCalls, assured.assuredCalls)
}

func TestWhenEndpointSuccess(t *testing.T) {
	assured := &AssuredEndpoints{
		assuredCalls: fullAssuredCalls,
		madeCalls:    map[string][]*AssuredCall{},
	}
	expected := map[string][]*AssuredCall{
		"test/assured":   []*AssuredCall{call2, call1},
		"teapot/assured": []*AssuredCall{call3},
	}

	c, err := assured.WhenEndpoint(ctx, call1)

	require.NoError(t, err)
	require.Equal(t, call1, c)
	require.Equal(t, expected, assured.assuredCalls)

	c, err = assured.WhenEndpoint(ctx, call2)

	require.NoError(t, err)
	require.Equal(t, call2, c)
	require.Equal(t, fullAssuredCalls, assured.assuredCalls)

	c, err = assured.WhenEndpoint(ctx, call3)

	require.NoError(t, err)
	require.Equal(t, call3, c)
	require.Equal(t, fullAssuredCalls, assured.assuredCalls)
	require.Equal(t, fullAssuredCalls, assured.madeCalls)
}

func TestWhenEndpointNotFound(t *testing.T) {
	assured := NewAssuredEndpoints(kitlog.NewLogfmtLogger(ioutil.Discard))

	c, err := assured.WhenEndpoint(ctx, call1)

	require.Nil(t, c)
	require.Error(t, err)
	require.Equal(t, "No assured calls", err.Error())
}

func TestThenEndpointSuccess(t *testing.T) {
	assured := &AssuredEndpoints{
		madeCalls: fullAssuredCalls,
	}

	c, err := assured.ThenEndpoint(ctx, call1)

	require.NoError(t, err)
	require.Equal(t, []*AssuredCall{call1, call2}, c)

	c, err = assured.ThenEndpoint(ctx, call3)

	require.NoError(t, err)
	require.Equal(t, []*AssuredCall{call3}, c)
}

func TestClearEndpointSuccess(t *testing.T) {
	assured := &AssuredEndpoints{
		logger:       kitlog.NewLogfmtLogger(ioutil.Discard),
		assuredCalls: fullAssuredCalls,
		madeCalls:    fullAssuredCalls,
	}
	expected := map[string][]*AssuredCall{
		"teapot/assured": []*AssuredCall{call3},
	}

	c, err := assured.ClearEndpoint(ctx, call1)

	require.NoError(t, err)
	require.Nil(t, c)
	require.Equal(t, expected, assured.assuredCalls)
	require.Equal(t, expected, assured.madeCalls)

	c, err = assured.ClearEndpoint(ctx, call2)

	require.NoError(t, err)
	require.Nil(t, c)
	require.Equal(t, expected, assured.assuredCalls)
	require.Equal(t, expected, assured.madeCalls)

	c, err = assured.ClearEndpoint(ctx, call3)

	require.NoError(t, err)
	require.Nil(t, c)
	require.Equal(t, map[string][]*AssuredCall{}, assured.assuredCalls)
	require.Equal(t, map[string][]*AssuredCall{}, assured.madeCalls)
}

func TestClearAllEndpointSuccess(t *testing.T) {
	assured := &AssuredEndpoints{
		logger:       kitlog.NewLogfmtLogger(ioutil.Discard),
		assuredCalls: fullAssuredCalls,
		madeCalls:    fullAssuredCalls,
	}

	c, err := assured.ClearAllEndpoint(ctx, nil)

	require.NoError(t, err)
	require.Nil(t, c)
	require.Equal(t, map[string][]*AssuredCall{}, assured.assuredCalls)
	require.Equal(t, map[string][]*AssuredCall{}, assured.madeCalls)
}

var (
	ctx   = context.Background()
	call1 = &AssuredCall{
		Path:       "test/assured",
		StatusCode: http.StatusOK,
		Response:   map[string]interface{}{"assured": true},
	}
	call2 = &AssuredCall{
		Path:       "test/assured",
		StatusCode: http.StatusConflict,
		Response:   "error",
	}
	call3 = &AssuredCall{
		Path:       "teapot/assured",
		StatusCode: http.StatusTeapot,
	}
	fullAssuredCalls = map[string][]*AssuredCall{
		"test/assured":   []*AssuredCall{call1, call2},
		"teapot/assured": []*AssuredCall{call3},
	}
)

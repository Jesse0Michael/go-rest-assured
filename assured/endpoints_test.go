package assured

import (
	"context"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	kitlog "github.com/go-kit/kit/log"
	"github.com/stretchr/testify/require"
)

func TestNewAssuredEndpoints(t *testing.T) {
	expected := &AssuredEndpoints{
		logger:         testSettings.Logger,
		httpClient:     *http.DefaultClient,
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
	c, err := actual(ctx, testCall1())

	require.NoError(t, err)
	require.Equal(t, testCall1(), c)
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

	c, err := endpoints.GivenEndpoint(ctx, testCall1())

	require.NoError(t, err)
	require.Equal(t, testCall1(), c)

	c, err = endpoints.GivenEndpoint(ctx, testCall2())

	require.NoError(t, err)
	require.Equal(t, testCall2(), c)

	c, err = endpoints.GivenEndpoint(ctx, testCall3())

	require.NoError(t, err)
	require.Equal(t, testCall3(), c)

	require.Equal(t, fullAssuredCalls, endpoints.assuredCalls)
}

func TestGivenCallbackEndpointSuccess(t *testing.T) {
	endpoints := NewAssuredEndpoints(testSettings)

	callback1 := testCall1()
	callback1.Headers[AssuredCallbackKey] = "call-key"
	c, err := endpoints.GivenEndpoint(ctx, callback1)

	require.NoError(t, err)
	require.Equal(t, callback1, c)

	callback2 := testCall2()
	callback2.Headers[AssuredCallbackKey] = "call-key"
	c, err = endpoints.GivenEndpoint(ctx, callback2)

	require.NoError(t, err)
	require.Equal(t, callback2, c)

	callback3 := testCall3()
	callback3.Headers[AssuredCallbackKey] = "call-key"
	c, err = endpoints.GivenEndpoint(ctx, callback3)

	require.NoError(t, err)
	require.Equal(t, callback3, c)

	c, err = endpoints.GivenCallbackEndpoint(ctx, testCallback())

	expectedAssured := &CallStore{
		data: map[string][]*Call{
			"GET:test/assured":    {callback1, callback2},
			"POST:teapot/assured": {callback3},
		},
	}
	expectedCallback := &CallStore{
		data: map[string][]*Call{
			"call-key": {testCallback()},
		},
	}
	require.NoError(t, err)
	require.Equal(t, testCallback(), c)
	require.Equal(t, expectedAssured, endpoints.assuredCalls)
	require.Equal(t, expectedCallback, endpoints.callbackCalls)

}

func TestWhenEndpointSuccess(t *testing.T) {
	endpoints := &AssuredEndpoints{
		logger:         testSettings.Logger,
		assuredCalls:   fullAssuredCalls,
		madeCalls:      NewCallStore(),
		callbackCalls:  NewCallStore(),
		trackMadeCalls: true,
	}
	expected := map[string][]*Call{
		"GET:test/assured":    {testCall2(), testCall1()},
		"POST:teapot/assured": {testCall3()},
	}

	c, err := endpoints.WhenEndpoint(ctx, testCall1())

	require.NoError(t, err)
	require.Equal(t, testCall1(), c)
	require.Equal(t, expected, endpoints.assuredCalls.data)

	c, err = endpoints.WhenEndpoint(ctx, testCall2())

	require.NoError(t, err)
	require.Equal(t, testCall2(), c)
	require.Equal(t, fullAssuredCalls, endpoints.assuredCalls)

	c, err = endpoints.WhenEndpoint(ctx, testCall3())

	require.NoError(t, err)
	require.Equal(t, testCall3(), c)
	require.Equal(t, fullAssuredCalls, endpoints.assuredCalls)
	require.Equal(t, fullAssuredCalls, endpoints.madeCalls)
}

func TestWhenEndpointSuccessTrackingDisabled(t *testing.T) {
	endpoints := &AssuredEndpoints{
		logger:         testSettings.Logger,
		assuredCalls:   fullAssuredCalls,
		madeCalls:      NewCallStore(),
		callbackCalls:  NewCallStore(),
		trackMadeCalls: false,
	}
	expected := map[string][]*Call{
		"GET:test/assured":    {testCall2(), testCall1()},
		"POST:teapot/assured": {testCall3()},
	}

	c, err := endpoints.WhenEndpoint(ctx, testCall1())

	require.NoError(t, err)
	require.Equal(t, testCall1(), c)
	require.Equal(t, expected, endpoints.assuredCalls.data)

	c, err = endpoints.WhenEndpoint(ctx, testCall2())

	require.NoError(t, err)
	require.Equal(t, testCall2(), c)
	require.Equal(t, fullAssuredCalls, endpoints.assuredCalls)

	c, err = endpoints.WhenEndpoint(ctx, testCall3())

	require.NoError(t, err)
	require.Equal(t, testCall3(), c)
	require.Equal(t, fullAssuredCalls, endpoints.assuredCalls)
	require.Equal(t, NewCallStore(), endpoints.madeCalls)
}

func TestWhenEndpointSuccessCallbacks(t *testing.T) {
	called := false
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	}))
	assured := testCall1()
	assured.Headers[AssuredCallbackKey] = "call-key"
	call := testCallback()
	call.Headers[AssuredCallbackTarget] = testServer.URL
	endpoints := &AssuredEndpoints{
		logger: testSettings.Logger,
		assuredCalls: &CallStore{
			data: map[string][]*Call{"GET:test/assured": []*Call{assured}},
		},
		madeCalls: NewCallStore(),
		callbackCalls: &CallStore{
			data: map[string][]*Call{"call-key": []*Call{call}},
		},
		trackMadeCalls: true,
	}

	c, err := endpoints.WhenEndpoint(ctx, assured)

	require.NoError(t, err)
	require.Equal(t, assured, c)
	// allow go routine to finish
	time.Sleep(1 * time.Millisecond)
	require.True(t, called, "callback was not hit")
}

func TestWhenEndpointSuccessCallbacksDelayed(t *testing.T) {
	called := false
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	}))
	assured := testCall1()
	assured.Headers[AssuredCallbackKey] = "call-key"
	call := testCallback()
	call.Headers[AssuredCallbackTarget] = testServer.URL
	call.Headers[AssuredCallbackDelay] = "2"
	endpoints := &AssuredEndpoints{
		logger: testSettings.Logger,
		assuredCalls: &CallStore{
			data: map[string][]*Call{"GET:test/assured": []*Call{assured}},
		},
		madeCalls: NewCallStore(),
		callbackCalls: &CallStore{
			data: map[string][]*Call{"call-key": []*Call{call}},
		},
		trackMadeCalls: true,
	}

	c, err := endpoints.WhenEndpoint(ctx, assured)

	require.NoError(t, err)
	require.Equal(t, assured, c)
	// allow go routine to finish
	time.Sleep(1 * time.Second)
	require.False(t, called, "callback should not be hit yet")
	time.Sleep(2 * time.Second)
	require.True(t, called, "callback was not hit")
}

func TestSendCallbackBadRequest(t *testing.T) {
	called := false
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	}))
	call := testCallback()
	call.Method = "\""
	endpoints := NewAssuredEndpoints(testSettings)
	endpoints.sendCallback(testServer.URL, call)

	// allow go routine to finish
	time.Sleep(1 * time.Millisecond)
	require.False(t, called, "callback should not be hit")
}

func TestSendCallbackBadResponse(t *testing.T) {
	endpoints := NewAssuredEndpoints(testSettings)
	endpoints.sendCallback("http://localhost:900000", testCallback())
}

func TestWhenEndpointNotFound(t *testing.T) {
	endpoints := NewAssuredEndpoints(testSettings)

	c, err := endpoints.WhenEndpoint(ctx, testCall1())

	require.Nil(t, c)
	require.Error(t, err)
	require.Equal(t, "No assured calls", err.Error())
}

func TestVerifyEndpointSuccess(t *testing.T) {
	endpoints := &AssuredEndpoints{
		madeCalls:      fullAssuredCalls,
		trackMadeCalls: true,
	}

	c, err := endpoints.VerifyEndpoint(ctx, testCall1())

	require.NoError(t, err)
	require.Equal(t, []*Call{testCall1(), testCall2()}, c)

	c, err = endpoints.VerifyEndpoint(ctx, testCall3())

	require.NoError(t, err)
	require.Equal(t, []*Call{testCall3()}, c)
}

func TestVerifyEndpointTrackingDisabled(t *testing.T) {
	endpoints := &AssuredEndpoints{
		madeCalls:      fullAssuredCalls,
		trackMadeCalls: false,
	}

	c, err := endpoints.VerifyEndpoint(ctx, testCall1())

	require.Nil(t, c)
	require.Error(t, err)
	require.Equal(t, "Tracking made calls is disabled", err.Error())
}

func TestClearEndpointSuccess(t *testing.T) {
	endpoints := &AssuredEndpoints{
		logger:         kitlog.NewLogfmtLogger(ioutil.Discard),
		assuredCalls:   fullAssuredCalls,
		madeCalls:      fullAssuredCalls,
		callbackCalls:  NewCallStore(),
		trackMadeCalls: true,
	}
	expected := map[string][]*Call{
		"POST:teapot/assured": {testCall3()},
	}

	c, err := endpoints.ClearEndpoint(ctx, testCall1())

	require.NoError(t, err)
	require.Nil(t, c)
	require.Equal(t, expected, endpoints.assuredCalls.data)
	require.Equal(t, expected, endpoints.madeCalls.data)

	c, err = endpoints.ClearEndpoint(ctx, testCall2())

	require.NoError(t, err)
	require.Nil(t, c)
	require.Equal(t, expected, endpoints.assuredCalls.data)
	require.Equal(t, expected, endpoints.madeCalls.data)

	c, err = endpoints.ClearEndpoint(ctx, testCall3())

	require.NoError(t, err)
	require.Nil(t, c)
	require.Equal(t, map[string][]*Call{}, endpoints.assuredCalls.data)
	require.Equal(t, map[string][]*Call{}, endpoints.madeCalls.data)
}

func TestClearEndpointSuccessCallback(t *testing.T) {
	endpoints := &AssuredEndpoints{
		logger:       kitlog.NewLogfmtLogger(ioutil.Discard),
		assuredCalls: fullAssuredCalls,
		madeCalls:    NewCallStore(),
		callbackCalls: &CallStore{
			data: map[string][]*Call{
				"call-key":       {testCallback()},
				"other-call-key": {testCallback()},
			},
		},
		trackMadeCalls: true,
	}

	c, err := endpoints.ClearEndpoint(ctx, testCallback())

	require.NoError(t, err)
	require.Nil(t, c)
	require.Equal(t, fullAssuredCalls.data, endpoints.assuredCalls.data)
	require.Equal(t, map[string][]*Call{}, endpoints.madeCalls.data)
	require.Equal(t, map[string][]*Call{"other-call-key": {testCallback()}}, endpoints.callbackCalls.data)
}

func TestClearAllEndpointSuccess(t *testing.T) {
	endpoints := &AssuredEndpoints{
		logger:         kitlog.NewLogfmtLogger(ioutil.Discard),
		assuredCalls:   fullAssuredCalls,
		madeCalls:      fullAssuredCalls,
		callbackCalls:  fullAssuredCalls,
		trackMadeCalls: true,
	}

	c, err := endpoints.ClearAllEndpoint(ctx, nil)

	require.NoError(t, err)
	require.Nil(t, c)
	require.Equal(t, map[string][]*Call{}, endpoints.assuredCalls.data)
	require.Equal(t, map[string][]*Call{}, endpoints.madeCalls.data)
	require.Equal(t, map[string][]*Call{}, endpoints.callbackCalls.data)
}

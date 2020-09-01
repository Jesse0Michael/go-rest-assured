package assured

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/require"
)

func TestApplicationRouterGivenBinding(t *testing.T) {
	router := NewClient().createApplicationRouter()

	for _, verb := range verbs {
		req, err := http.NewRequest(verb, "/given/rest/assured", nil)
		require.NoError(t, err)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)
		require.Equal(t, http.StatusOK, resp.Code)
		require.Equal(t, "*", resp.Header().Get("Access-Control-Allow-Origin"))
	}
}

func TestApplicationRouterGivenCallbackBinding(t *testing.T) {
	router := NewClient().createApplicationRouter()

	for _, verb := range verbs {
		req, err := http.NewRequest(verb, "/callback", nil)
		req.Header.Set(AssuredCallbackKey, "call-key")
		req.Header.Set(AssuredCallbackTarget, "http://faketarget.com/")
		require.NoError(t, err)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)
		require.Equal(t, http.StatusCreated, resp.Code)
		require.Equal(t, "*", resp.Header().Get("Access-Control-Allow-Origin"))
	}
}

func TestApplicationRouterWhenBinding(t *testing.T) {
	router := NewClient().createApplicationRouter()

	for _, verb := range verbs {
		req, err := http.NewRequest(verb, "/given/rest/assured", bytes.NewBuffer([]byte(`{"assured": true}`)))
		require.NoError(t, err)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		req, err = http.NewRequest(verb, "/when/rest/assured", nil)
		require.NoError(t, err)
		resp = httptest.NewRecorder()
		router.ServeHTTP(resp, req)
		require.Equal(t, http.StatusOK, resp.Code)
		require.Equal(t, "*", resp.Header().Get("Access-Control-Allow-Origin"))
	}
}

func TestApplicationRouterVerifyBinding(t *testing.T) {
	router := NewClient().createApplicationRouter()

	for _, verb := range verbs {
		req, err := http.NewRequest(verb, "/verify/rest/assured", nil)
		require.NoError(t, err)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)
		require.Equal(t, http.StatusOK, resp.Code)
		require.Equal(t, "*", resp.Header().Get("Access-Control-Allow-Origin"))
	}
}

func TestApplicationRouterClearBinding(t *testing.T) {
	router := NewClient().createApplicationRouter()

	for _, verb := range verbs {
		req, err := http.NewRequest(verb, "/clear/rest/assured", nil)
		require.NoError(t, err)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)
		require.Equal(t, http.StatusOK, resp.Code)
		require.Equal(t, "*", resp.Header().Get("Access-Control-Allow-Origin"))
	}
}

func TestApplicationRouterClearAllBinding(t *testing.T) {
	router := NewClient().createApplicationRouter()

	req, err := http.NewRequest(http.MethodDelete, "/clear", nil)
	require.NoError(t, err)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	require.Equal(t, http.StatusOK, resp.Code)
	require.Equal(t, "*", resp.Header().Get("Access-Control-Allow-Origin"))
}

func TestApplicationRouterFailure(t *testing.T) {
	router := NewClient().createApplicationRouter()

	req, err := http.NewRequest(http.MethodGet, "/trouble", nil)
	require.NoError(t, err)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	require.Equal(t, http.StatusNotFound, resp.Code)
}

func TestDecodeAssuredCall(t *testing.T) {
	decoded := false
	expected := &Call{
		Path:       "test/assured",
		StatusCode: http.StatusOK,
		Method:     http.MethodPost,
		Response:   []byte(`{"assured": true}`),
		Headers:    map[string]string{},
		Query:      map[string]string{"assured": "max"},
	}
	testDecode := func(resp http.ResponseWriter, req *http.Request) {
		c, err := decodeAssuredCall(context.TODO(), req)

		require.NoError(t, err)
		require.Equal(t, expected, c)
		decoded = true
	}

	req, err := http.NewRequest(http.MethodPost, "/verify/test/assured?assured=max", bytes.NewBuffer([]byte(`{"assured": true}`)))
	require.NoError(t, err)

	router := mux.NewRouter()
	router.HandleFunc("/verify/{path:.*}", testDecode).Methods(http.MethodPost)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	require.True(t, decoded, "decode method was not hit")
}

func TestDecodeAssuredCallNilBody(t *testing.T) {
	decoded := false
	expected := &Call{
		Path:       "test/assured",
		StatusCode: http.StatusOK,
		Method:     http.MethodDelete,
		Headers:    map[string]string{},
		Query:      map[string]string{},
	}
	testDecode := func(resp http.ResponseWriter, req *http.Request) {
		c, err := decodeAssuredCall(context.TODO(), req)

		require.NoError(t, err)
		require.Equal(t, expected, c)
		decoded = true
	}

	req, err := http.NewRequest(http.MethodDelete, "/when/test/assured", nil)
	require.NoError(t, err)

	router := mux.NewRouter()
	router.HandleFunc("/when/{path:.*}", testDecode).Methods(http.MethodDelete)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	require.True(t, decoded, "decode method was not hit")
}

func TestDecodeAssuredCallStatus(t *testing.T) {
	decoded := false
	expected := &Call{
		Path:       "test/assured",
		StatusCode: http.StatusForbidden,
		Method:     http.MethodGet,
		Headers:    map[string]string{"Assured-Status": "403"},
		Query:      map[string]string{},
	}
	testDecode := func(resp http.ResponseWriter, req *http.Request) {
		c, err := decodeAssuredCall(context.TODO(), req)

		require.NoError(t, err)
		require.Equal(t, expected, c)
		decoded = true
	}

	req, err := http.NewRequest(http.MethodGet, "/given/test/assured", nil)
	require.NoError(t, err)
	req.Header.Set("Assured-Status", "403")

	router := mux.NewRouter()
	router.HandleFunc("/given/{path:.*}", testDecode).Methods(http.MethodGet)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	require.True(t, decoded, "decode method was not hit")
}

func TestDecodeAssuredCallStatusFailure(t *testing.T) {
	decoded := false
	expected := &Call{
		Path:       "test/assured",
		StatusCode: http.StatusOK,
		Method:     http.MethodGet,
		Headers:    map[string]string{"Assured-Status": "four oh three"},
		Query:      map[string]string{},
	}
	testDecode := func(resp http.ResponseWriter, req *http.Request) {
		c, err := decodeAssuredCall(context.TODO(), req)

		require.NoError(t, err)
		require.Equal(t, expected, c)
		decoded = true
	}

	req, err := http.NewRequest(http.MethodGet, "/given/test/assured", nil)
	require.NoError(t, err)
	req.Header.Set("Assured-Status", "four oh three")

	router := mux.NewRouter()
	router.HandleFunc("/given/{path:.*}", testDecode).Methods(http.MethodGet)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	require.True(t, decoded, "decode method was not hit")
}

func TestDecodeAssuredCallback(t *testing.T) {
	decoded := false
	expected := &Call{
		Method:     http.MethodPost,
		StatusCode: http.StatusCreated,
		Response:   []byte(`{"done": true}`),
		Headers:    map[string]string{"Assured-Callback-Target": "http://faketarget.com/", "Assured-Callback-Key": "call-key"},
	}
	testDecode := func(resp http.ResponseWriter, req *http.Request) {
		c, err := decodeAssuredCallback(context.TODO(), req)

		require.NoError(t, err)
		require.Equal(t, expected, c)
		decoded = true
	}

	req, err := http.NewRequest(http.MethodPost, "/callback?assured=max", bytes.NewBuffer([]byte(`{"done": true}`)))
	require.NoError(t, err)
	req.Header.Set(AssuredCallbackKey, "call-key")
	req.Header.Set(AssuredCallbackTarget, "http://faketarget.com/")

	router := mux.NewRouter()
	router.HandleFunc("/callback", testDecode).Methods(http.MethodPost)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	require.True(t, decoded, "decode method was not hit")
}

func TestDecodeAssuredCallbackMissingKey(t *testing.T) {
	decoded := false
	testDecode := func(resp http.ResponseWriter, req *http.Request) {
		c, err := decodeAssuredCallback(context.TODO(), req)

		require.Nil(t, c)
		require.Error(t, err)
		require.Equal(t, "'Assured-Callback-Key' header required for callback", err.Error())
		decoded = true
	}

	req, err := http.NewRequest(http.MethodPost, "/callback", bytes.NewBuffer([]byte(`{"done": true}`)))
	req.Header.Set(AssuredCallbackTarget, "http://faketarget.com/")
	require.NoError(t, err)

	router := mux.NewRouter()
	router.HandleFunc("/callback", testDecode).Methods(http.MethodPost)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	require.True(t, decoded, "decode method was not hit")
}

func TestDecodeAssuredCallbackMissingTarget(t *testing.T) {
	decoded := false
	testDecode := func(resp http.ResponseWriter, req *http.Request) {
		c, err := decodeAssuredCallback(context.TODO(), req)

		require.Nil(t, c)
		require.Error(t, err)
		require.Equal(t, "'Assured-Callback-Target' header required for callback", err.Error())
		decoded = true
	}

	req, err := http.NewRequest(http.MethodPost, "/callback", bytes.NewBuffer([]byte(`{"done": true}`)))
	req.Header.Set(AssuredCallbackKey, "call-key")
	require.NoError(t, err)

	router := mux.NewRouter()
	router.HandleFunc("/callback", testDecode).Methods(http.MethodPost)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	require.True(t, decoded, "decode method was not hit")
}

func TestEncodeAssuredCall(t *testing.T) {
	call := &Call{
		Path:       "/test/assured",
		StatusCode: http.StatusCreated,
		Method:     http.MethodPost,
		Response:   []byte(`{"assured": true}`),
		Headers:    map[string]string{"Content-Length": "19", "User-Agent": "Go-http-client/1.1", "Accept-Encoding": "gzip", "Assured-Status": "403"},
	}
	resp := httptest.NewRecorder()

	err := encodeAssuredCall(context.TODO(), resp, call)

	require.NoError(t, err)
	require.Equal(t, http.StatusCreated, resp.Code)
	require.Equal(t, `{"assured": true}`, resp.Body.String())
	require.Equal(t, "19", resp.Header().Get("Content-Length"))
	require.Equal(t, "Go-http-client/1.1", resp.Header().Get("User-Agent"))
	require.Equal(t, "gzip", resp.Header().Get("Accept-Encoding"))
	require.Empty(t, resp.Header().Get("Assured-Status"))
}

func TestEncodeAssuredCalls(t *testing.T) {
	resp := httptest.NewRecorder()
	expected, err := ioutil.ReadFile("testdata/calls.json")
	require.NoError(t, err)
	err = encodeAssuredCall(context.TODO(), resp, []*Call{testCall1(), testCall2(), testCall3()})

	require.NoError(t, err)
	require.Equal(t, "application/json", resp.Result().Header.Get("Content-Type"))
	require.JSONEq(t, string(expected), resp.Body.String())
}

//go-rest-assured test vars
var (
	verbs = []string{
		http.MethodGet,
		http.MethodHead,
		http.MethodPost,
		http.MethodPut,
		http.MethodPatch,
		http.MethodDelete,
		http.MethodConnect,
		http.MethodOptions,
	}
	fullAssuredCalls = &CallStore{
		data: map[string][]*Call{
			"GET:test/assured":    {testCall1(), testCall2()},
			"POST:teapot/assured": {testCall3()},
		},
	}
)

func testCall1() *Call {
	return &Call{
		Path:       "test/assured",
		Method:     "GET",
		StatusCode: http.StatusOK,
		Response:   []byte(`{"assured": true}`),
		Headers:    map[string]string{"Content-Length": "17", "User-Agent": "Go-http-client/1.1", "Accept-Encoding": "gzip"},
		Query:      map[string]string{"assured": "max"},
	}
}

func testCall2() *Call {
	return &Call{
		Path:       "test/assured",
		Method:     "GET",
		StatusCode: http.StatusConflict,
		Response:   []byte("error"),
		Headers:    map[string]string{"Content-Length": "5", "User-Agent": "Go-http-client/1.1", "Accept-Encoding": "gzip"},
	}
}

func testCall3() *Call {
	return &Call{
		Path:       "teapot/assured",
		Method:     "POST",
		StatusCode: http.StatusTeapot,
		Headers:    map[string]string{"Content-Length": "0", "User-Agent": "Go-http-client/1.1", "Accept-Encoding": "gzip"},
	}
}

func testCallback() *Call {
	return &Call{
		Response: []byte(`{"done": true}`),
		Method:   "POST",
		Headers:  map[string]string{"Assured-Callback-Key": "call-key", "Assured-Callback-Target": "http://faketarget.com/"},
	}
}

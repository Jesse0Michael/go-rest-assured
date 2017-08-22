package bindings

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	kitlog "github.com/go-kit/kit/log"
	"github.com/jesse0michael/go-rest-assured/assured"
	"github.com/stretchr/testify/require"
)

func TestApplicationRouterGivenBinding(t *testing.T) {
	router := createApplicationRouter(ctx, kitlog.NewLogfmtLogger(ioutil.Discard))

	for _, verb := range verbs {
		req, err := http.NewRequest(verb, "/given/rest/assured", nil)
		require.NoError(t, err)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)
		require.Equal(t, http.StatusOK, resp.Code)
		require.Equal(t, "*", resp.Header().Get("Access-Control-Allow-Origin"))
	}
}

func TestApplicationRouterWhenBinding(t *testing.T) {
	router := createApplicationRouter(ctx, kitlog.NewLogfmtLogger(ioutil.Discard))

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

func TestApplicationRouterThenBinding(t *testing.T) {
	router := createApplicationRouter(ctx, kitlog.NewLogfmtLogger(ioutil.Discard))

	for _, verb := range verbs {
		req, err := http.NewRequest(verb, "/then/rest/assured", nil)
		require.NoError(t, err)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)
		require.Equal(t, http.StatusOK, resp.Code)
		require.Equal(t, "*", resp.Header().Get("Access-Control-Allow-Origin"))
	}
}

func TestApplicationRouterClearBinding(t *testing.T) {
	router := createApplicationRouter(ctx, kitlog.NewLogfmtLogger(ioutil.Discard))

	req, err := http.NewRequest(http.MethodDelete, "/clear/rest/assured", nil)
	require.NoError(t, err)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	require.Equal(t, http.StatusOK, resp.Code)
	require.Equal(t, "*", resp.Header().Get("Access-Control-Allow-Origin"))
}

func TestApplicationRouterClearAllBinding(t *testing.T) {
	router := createApplicationRouter(ctx, kitlog.NewLogfmtLogger(ioutil.Discard))

	req, err := http.NewRequest(http.MethodDelete, "/clear", nil)
	require.NoError(t, err)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	require.Equal(t, http.StatusOK, resp.Code)
	require.Equal(t, "*", resp.Header().Get("Access-Control-Allow-Origin"))
}

func TestApplicationRouterFailure(t *testing.T) {
	router := createApplicationRouter(ctx, kitlog.NewLogfmtLogger(ioutil.Discard))

	req, err := http.NewRequest(http.MethodGet, "/trouble", nil)
	require.NoError(t, err)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	require.Equal(t, http.StatusNotFound, resp.Code)
}

func TestDecodeAssuredCall(t *testing.T) {
	req, err := http.NewRequest(http.MethodPost, "/then/test/assured?test=positive", bytes.NewBuffer([]byte(`{"assured": true}`)))
	require.NoError(t, err)
	expected := &assured.Call{
		Path:       "/then/test/assured",
		StatusCode: http.StatusOK,
		Method:     http.MethodPost,
		Response:   []byte(`{"assured": true}`),
	}

	c, err := decodeAssuredCall(ctx, req)

	require.NoError(t, err)
	require.Equal(t, expected, c)
}

func TestDecodeAssuredCallNilBody(t *testing.T) {
	req, err := http.NewRequest(http.MethodDelete, "/when/test/assured", nil)
	require.NoError(t, err)
	expected := &assured.Call{
		Path:       "/when/test/assured",
		StatusCode: http.StatusOK,
		Method:     http.MethodDelete,
	}

	c, err := decodeAssuredCall(ctx, req)

	require.NoError(t, err)
	require.Equal(t, expected, c)
}

func TestDecodeAssuredCallStatus(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "/given/test/assured", nil)
	require.NoError(t, err)
	req.Header.Set("Assured-Status", "403")
	expected := &assured.Call{
		Path:       "/given/test/assured",
		StatusCode: http.StatusForbidden,
		Method:     http.MethodGet,
	}

	c, err := decodeAssuredCall(ctx, req)

	require.NoError(t, err)
	require.Equal(t, expected, c)
}

func TestDecodeAssuredCallStatusFailure(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "/given/test/assured", nil)
	require.NoError(t, err)
	req.Header.Set("Assured-Status", "four oh three")
	expected := &assured.Call{
		Path:       "/given/test/assured",
		StatusCode: http.StatusOK,
		Method:     http.MethodGet,
	}

	c, err := decodeAssuredCall(ctx, req)

	require.NoError(t, err)
	require.Equal(t, expected, c)
}

func TestEncodeAssuredCall(t *testing.T) {
	call := &assured.Call{
		Path:       "/then/test/assured",
		StatusCode: http.StatusCreated,
		Method:     http.MethodPost,
		Response:   []byte(`{"assured": true}`),
	}
	resp := httptest.NewRecorder()

	err := encodeAssuredCall(ctx, resp, call)

	require.NoError(t, err)
	require.Equal(t, http.StatusCreated, resp.Code)
	require.Equal(t, `{"assured": true}`, resp.Body.String())
}

var (
	ctx   = context.Background()
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
	paths = []string{
		"/given/test/assured",
		"/when/test/assured",
		"/then/test/assured",
	}
)

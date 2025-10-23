package assured

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestAssured(t *testing.T) {
	assured := NewAssured(WithPort(9091))
	go func() { _ = assured.Serve() }()
	defer func() { _ = assured.Close() }()
	time.Sleep(time.Second)

	url := assured.URL()
	require.Equal(t, "http://localhost:9091", url)
	require.NoError(t, assured.Given(t.Context(), *testCall1()))
	require.NoError(t, assured.Given(t.Context(), *testCall2()))
	require.NoError(t, assured.Given(t.Context(), *testCall3()))

	req, err := http.NewRequest(http.MethodGet, url+"/test/assured", bytes.NewReader([]byte(`{"calling":"you"}`)))
	require.NoError(t, err)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	require.Equal(t, []byte(`{"assured": true}`), body)

	req, err = http.NewRequest(http.MethodGet, url+"/test/assured", bytes.NewReader([]byte(`{"calling":"again"}`)))
	require.NoError(t, err)

	resp, err = http.DefaultClient.Do(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusConflict, resp.StatusCode)
	body, err = io.ReadAll(resp.Body)
	require.NoError(t, err)
	require.Equal(t, []byte("error"), body)

	req, err = http.NewRequest(http.MethodPost, url+"/teapot/assured", bytes.NewReader([]byte(`{"calling":"here"}`)))
	require.NoError(t, err)

	resp, err = http.DefaultClient.Do(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusTeapot, resp.StatusCode)
	body, err = io.ReadAll(resp.Body)
	require.NoError(t, err)
	require.Equal(t, []byte{}, body)

	calls, err := assured.Verify(t.Context(), http.MethodGet, "test/assured")
	require.NoError(t, err)
	require.Equal(t, []Call{
		{
			Method:   http.MethodGet,
			Path:     "test/assured",
			Response: []byte(`{"calling":"you"}`),
			Headers:  map[string]string{"Content-Length": "17", "User-Agent": "Go-http-client/1.1", "Accept-Encoding": "gzip"}},
		{
			Method:   http.MethodGet,
			Path:     "test/assured",
			Response: []byte(`{"calling":"again"}`),
			Headers:  map[string]string{"Content-Length": "19", "User-Agent": "Go-http-client/1.1", "Accept-Encoding": "gzip"}}}, calls)

	calls, err = assured.Verify(t.Context(), http.MethodPost, "teapot/assured")
	require.NoError(t, err)
	require.Equal(t, []Call{
		{
			Method:   http.MethodPost,
			Path:     "teapot/assured",
			Response: []byte(`{"calling":"here"}`),
			Headers:  map[string]string{"Content-Length": "18", "User-Agent": "Go-http-client/1.1", "Accept-Encoding": "gzip"}}}, calls)

	err = assured.Clear(t.Context(), http.MethodGet, "test/assured")
	require.NoError(t, err)

	calls, err = assured.Verify(t.Context(), http.MethodGet, "test/assured")
	require.NoError(t, err)
	require.Nil(t, calls)

	calls, err = assured.Verify(t.Context(), http.MethodPost, "teapot/assured")
	require.NoError(t, err)
	require.Equal(t, []Call{
		{
			Method:   http.MethodPost,
			Path:     "teapot/assured",
			Response: []byte(`{"calling":"here"}`),
			Headers:  map[string]string{"Content-Length": "18", "User-Agent": "Go-http-client/1.1", "Accept-Encoding": "gzip"},
		},
	}, calls)

	err = assured.ClearAll(t.Context())
	require.NoError(t, err)

	calls, err = assured.Verify(t.Context(), http.MethodGet, "test/assured")
	require.NoError(t, err)
	require.Nil(t, calls)

	calls, err = assured.Verify(t.Context(), http.MethodPost, "teapot/assured")
	require.NoError(t, err)
	require.Nil(t, calls)
}

func TestAssuredTLS(t *testing.T) {
	insecureClient := http.Client{Transport: &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}}
	assured := NewAssured(WithTLS("testdata/localhost.pem", "testdata/localhost-key.pem"), WithPort(9092),
		WithHTTPClient(insecureClient))
	go func() { _ = assured.Serve() }()
	defer func() { _ = assured.Close() }()
	time.Sleep(1 * time.Second)

	url := assured.URL()
	require.Equal(t, "https://localhost:9092", url)
	require.NoError(t, assured.Given(t.Context(), *testCall1()))

	req, err := http.NewRequest(http.MethodGet, url+"/test/assured", bytes.NewReader([]byte(`{"calling":"you"}`)))
	require.NoError(t, err)

	resp, err := insecureClient.Do(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	require.Equal(t, []byte(`{"assured": true}`), body)

	calls, err := assured.Verify(t.Context(), http.MethodGet, "test/assured")
	require.NoError(t, err)
	require.Equal(t, []Call{
		{
			Method:   http.MethodGet,
			Path:     "test/assured",
			Response: []byte(`{"calling":"you"}`),
			Headers:  map[string]string{"Content-Length": "17", "User-Agent": "Go-http-client/1.1", "Accept-Encoding": "gzip"},
		},
	}, calls)
}

func TestAssuredCallbacks(t *testing.T) {
	httpClient := http.Client{}
	called := false
	delayCalled := false
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		require.NoError(t, err)
		require.Equal(t, []byte(`{"done":"here"}`), body)
		require.NotEmpty(t, r.Header.Get("x-info"))
		called = true
	}))
	delayTestServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		require.NoError(t, err)
		require.Equal(t, []byte(`{"wait":"there's more"}`), body)
		delayCalled = true
	}))
	assured := NewAssured()
	go func() { _ = assured.Serve() }()
	defer func() { _ = assured.Close() }()
	time.Sleep(time.Second)

	err := assured.Given(t.Context(), Call{
		Path:   "test/assured",
		Method: http.MethodPost,
		Delay:  2,
		Callbacks: []Callback{
			{
				Method:   http.MethodPost,
				Target:   testServer.URL,
				Response: []byte(`{"done":"here"}`),
				Headers:  map[string]string{"x-info": "important"},
			},
			{
				Method:   http.MethodPost,
				Target:   delayTestServer.URL,
				Delay:    4,
				Response: []byte(`{"wait":"there's more"}`),
			},
		},
	})
	require.NoError(t, err)

	req, err := http.NewRequest(http.MethodPost, assured.URL()+"/test/assured", bytes.NewReader([]byte(`{"calling":"here"}`)))
	require.NoError(t, err)

	start := time.Now()
	_, err = httpClient.Do(req)
	require.NoError(t, err)

	require.True(t, time.Since(start) >= 2*time.Second, "response should be delayed 2 seconds")
	// allow go routine to finish
	time.Sleep(1 * time.Second)
	require.True(t, called, "callback was not hit")
	require.False(t, delayCalled, "delayed callback should not be hit yet")
	time.Sleep(2 * time.Second)
	require.True(t, delayCalled, "delayed callback was not hit")
}

func TestAssuredClose(t *testing.T) {
	assured := NewAssured()
	go func() { _ = assured.Serve() }()
	assured2 := NewAssured()
	go func() { _ = assured2.Serve() }()
	time.Sleep(time.Second)

	require.NotEqual(t, assured.URL(), assured2.URL())

	require.NoError(t, assured.Given(t.Context(), *testCall1()))
	require.NoError(t, assured2.Given(t.Context(), *testCall1()))

	err := assured.Close()
	require.NoError(t, err)
	time.Sleep(time.Second)
	err = assured.Given(t.Context(), *testCall1())

	require.Error(t, err)
	require.Contains(t, err.Error(), `connection refused`)

	err = assured2.Close()
	require.NoError(t, err)
	time.Sleep(time.Second)
	err = assured2.Given(t.Context(), *testCall1())

	require.Error(t, err)
	require.Contains(t, err.Error(), `connection refused`)
}

func TestAssuredGivenNoMethod(t *testing.T) {
	assured := NewAssured()
	go func() { _ = assured.Serve() }()
	defer func() { _ = assured.Close() }()
	time.Sleep(time.Second)

	err := assured.Given(t.Context(), Call{Path: "NoMethodMan"})
	require.NoError(t, err)

	req, err := http.NewRequest(http.MethodGet, assured.URL()+"/NoMethodMan", nil)
	require.NoError(t, err)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestAssuredGivenCallbackMissingTarget(t *testing.T) {
	call := Call{
		Method: http.MethodPost,
		Callbacks: []Callback{
			{Method: http.MethodPost},
		},
	}
	assured := NewAssured()
	go func() { _ = assured.Serve() }()
	defer func() { _ = assured.Close() }()

	err := assured.Given(t.Context(), call)

	require.Error(t, err)
	require.Equal(t, "400:cannot stub callback without target", err.Error())
}

func TestAssuredGivenCallbackBadMethod(t *testing.T) {
	call := Call{
		Method: http.MethodPost,
		Callbacks: []Callback{
			{Method: "\"", Target: "http://localhost/"},
		},
	}
	assured := NewAssured()
	go func() { _ = assured.Serve() }()
	defer func() { _ = assured.Close() }()

	err := assured.Given(t.Context(), call)

	require.Error(t, err)
	require.Equal(t, "400:net/http: invalid method \"\\\"\"", err.Error())
}

func TestAssuredBadRequestFailure(t *testing.T) {
	assured := NewAssured()
	go func() { _ = assured.Serve() }()
	defer func() { _ = assured.Close() }()

	err := assured.Given(t.Context(), Call{Method: "\"", Path: "goat/path"})

	require.Error(t, err)
	require.Equal(t, `net/http: invalid method "\""`, err.Error())

	err = assured.Given(t.Context(), Call{Method: "\"", Path: "goat/path", Response: []byte("goats among men")})

	require.Error(t, err)
	require.Equal(t, `net/http: invalid method "\""`, err.Error())

	calls, err := assured.Verify(t.Context(), "\"", "goat/path")

	require.Error(t, err)
	require.Equal(t, `400:net/http: invalid method "\""`, err.Error())
	require.Nil(t, calls)

	err = assured.Clear(t.Context(), "\"", "goat/path")

	require.Error(t, err)
	require.Equal(t, `400:net/http: invalid method "\""`, err.Error())

	assured.baseURL = "http://localhost:-1"
	err = assured.ClearAll(t.Context())

	require.Error(t, err)
	require.Equal(t, `parse "http://localhost:-1/assured/clearall": invalid port ":-1" after host`, err.Error())
}

func TestAssuredVerifyHttpClientFailure(t *testing.T) {
	assured := NewAssured()
	go func() { _ = assured.Serve() }()
	err := assured.Close()
	require.NoError(t, err)

	calls, err := assured.Verify(t.Context(), "GONE", "not/started")

	require.Error(t, err)
	require.Contains(t, err.Error(), `connection refused`)
	require.Nil(t, calls)
}

func TestAssuredVerifyBodyFailure(t *testing.T) {
	assured := NewAssured()
	go func() { _ = assured.Serve() }()
	defer func() { _ = assured.Close() }()
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode("ydob+dab")
	}))
	defer testServer.Close()
	index := strings.LastIndex(testServer.URL, ":")
	port, err := strconv.ParseInt(testServer.URL[index+1:], 10, 64)
	require.NoError(t, err)
	assured.baseURL = fmt.Sprintf("http://localhost:%d", port)

	calls, err := assured.Verify(t.Context(), "BODY", "bad+body")

	require.Error(t, err)
	require.Equal(t, `json: cannot unmarshal string into Go value of type []assured.Call`, err.Error())
	require.Nil(t, calls)
}

func TestAssuredPathSanitization(t *testing.T) {
	assured := NewAssured()
	go func() { _ = assured.Serve() }()
	defer func() { _ = assured.Close() }()
	time.Sleep(time.Second)

	require.NoError(t, assured.Given(t.Context(), Call{Method: http.MethodGet, Path: "///yoyo/path///", StatusCode: http.StatusAccepted}))

	req, err := http.NewRequest(http.MethodGet, assured.URL()+"/yoyo/path", nil)
	require.NoError(t, err)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusAccepted, resp.StatusCode)
}

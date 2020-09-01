package assured

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestClient(t *testing.T) {
	client := NewClient(WithPort(9091))
	time.Sleep(time.Second)

	url := client.URL()
	require.Equal(t, "http://localhost:9091/when", url)
	require.NoError(t, client.Given(*testCall1()))
	require.NoError(t, client.Given(*testCall2()))
	require.NoError(t, client.Given(*testCall3()))

	req, err := http.NewRequest(http.MethodGet, url+"/test/assured", bytes.NewReader([]byte(`{"calling":"you"}`)))
	require.NoError(t, err)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	body, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)
	require.Equal(t, []byte(`{"assured": true}`), body)

	req, err = http.NewRequest(http.MethodGet, url+"/test/assured", bytes.NewReader([]byte(`{"calling":"again"}`)))
	require.NoError(t, err)

	resp, err = http.DefaultClient.Do(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusConflict, resp.StatusCode)
	body, err = ioutil.ReadAll(resp.Body)
	require.NoError(t, err)
	require.Equal(t, []byte("error"), body)

	req, err = http.NewRequest(http.MethodPost, url+"/teapot/assured", bytes.NewReader([]byte(`{"calling":"here"}`)))
	require.NoError(t, err)

	resp, err = http.DefaultClient.Do(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusTeapot, resp.StatusCode)
	body, err = ioutil.ReadAll(resp.Body)
	require.NoError(t, err)
	require.Equal(t, []byte{}, body)

	calls, err := client.Verify("GET", "test/assured")
	require.NoError(t, err)
	require.Equal(t, []Call{
		{
			Method:     "GET",
			Path:       "test/assured",
			StatusCode: 200,
			Response:   []byte(`{"calling":"you"}`),
			Headers:    map[string]string{"Content-Length": "17", "User-Agent": "Go-http-client/1.1", "Accept-Encoding": "gzip"}},
		{
			Method:     "GET",
			Path:       "test/assured",
			StatusCode: 200,
			Response:   []byte(`{"calling":"again"}`),
			Headers:    map[string]string{"Content-Length": "19", "User-Agent": "Go-http-client/1.1", "Accept-Encoding": "gzip"}}}, calls)

	calls, err = client.Verify("POST", "teapot/assured")
	require.NoError(t, err)
	require.Equal(t, []Call{
		{
			Method:     "POST",
			Path:       "teapot/assured",
			StatusCode: 200,
			Response:   []byte(`{"calling":"here"}`),
			Headers:    map[string]string{"Content-Length": "18", "User-Agent": "Go-http-client/1.1", "Accept-Encoding": "gzip"}}}, calls)

	err = client.Clear("GET", "test/assured")
	require.NoError(t, err)

	calls, err = client.Verify("GET", "test/assured")
	require.NoError(t, err)
	require.Nil(t, calls)

	calls, err = client.Verify("POST", "teapot/assured")
	require.NoError(t, err)
	require.Equal(t, []Call{
		{
			Method:     "POST",
			Path:       "teapot/assured",
			StatusCode: 200,
			Response:   []byte(`{"calling":"here"}`),
			Headers:    map[string]string{"Content-Length": "18", "User-Agent": "Go-http-client/1.1", "Accept-Encoding": "gzip"}}}, calls)

	err = client.ClearAll()
	require.NoError(t, err)

	calls, err = client.Verify("GET", "test/assured")
	require.NoError(t, err)
	require.Nil(t, calls)

	calls, err = client.Verify("POST", "teapot/assured")
	require.NoError(t, err)
	require.Nil(t, calls)
}

func TestClientCallbacks(t *testing.T) {
	httpClient := http.Client{}
	called := false
	delayCalled := false
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		require.NoError(t, err)
		require.Equal(t, []byte(`{"done":"here"}`), body)
		require.NotEmpty(t, r.Header.Get("x-info"))
		called = true
	}))
	delayTestServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		require.NoError(t, err)
		require.Equal(t, []byte(`{"wait":"there's more"}`), body)
		delayCalled = true
	}))
	client := NewClient()
	time.Sleep(time.Second)

	require.NoError(t, client.Given(Call{
		Path:   "test/assured",
		Method: "POST",
		Delay:  2,
		Callbacks: []Callback{
			{
				Method:   "POST",
				Target:   testServer.URL,
				Response: []byte(`{"done":"here"}`),
				Headers:  map[string]string{"x-info": "important"},
			},
			{
				Method:   "POST",
				Target:   delayTestServer.URL,
				Delay:    4,
				Response: []byte(`{"wait":"there's more"}`),
			},
		},
	}))

	req, err := http.NewRequest(http.MethodPost, client.URL()+"/test/assured", bytes.NewReader([]byte(`{"calling":"here"}`)))
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

func TestClientClose(t *testing.T) {
	client := NewClient()
	client2 := NewClient()
	time.Sleep(time.Second)

	require.NotEqual(t, client.URL(), client2.URL())

	require.NoError(t, client.Given(*testCall1()))
	require.NoError(t, client2.Given(*testCall1()))

	client.Close()
	time.Sleep(time.Second)
	err := client.Given(*testCall1())

	require.Error(t, err)
	require.Contains(t, err.Error(), `connection refused`)

	client2.Close()
	time.Sleep(time.Second)
	err = client2.Given(*testCall1())

	require.Error(t, err)
	require.Contains(t, err.Error(), `connection refused`)
}

func TestClientGivenNoMethod(t *testing.T) {
	client := NewClient()
	time.Sleep(time.Second)

	err := client.Given(Call{Path: "NoMethodMan"})
	require.NoError(t, err)

	req, err := http.NewRequest(http.MethodGet, client.URL()+"/NoMethodMan", nil)
	require.NoError(t, err)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestClientGivenCallbackMissingTarget(t *testing.T) {
	call := Call{
		Method: "POST",
		Callbacks: []Callback{
			{Method: "POST"},
		},
	}
	client := NewClient()

	err := client.Given(call)

	require.Error(t, err)
	require.Equal(t, "cannot stub callback without target", err.Error())
}

func TestClientGivenCallbackBadMethod(t *testing.T) {
	call := Call{
		Method: "POST",
		Callbacks: []Callback{
			{Method: "\"", Target: "http://localhost/"},
		},
	}
	client := NewClient()

	err := client.Given(call)

	require.Error(t, err)
	require.Equal(t, "net/http: invalid method \"\\\"\"", err.Error())
}

func TestClientBadRequestFailure(t *testing.T) {
	client := NewClient()

	err := client.Given(Call{Method: "\"", Path: "goat/path"})

	require.Error(t, err)
	require.Equal(t, `net/http: invalid method "\""`, err.Error())

	err = client.Given(Call{Method: "\"", Path: "goat/path", Response: []byte("goats among men")})

	require.Error(t, err)
	require.Equal(t, `net/http: invalid method "\""`, err.Error())

	calls, err := client.Verify("\"", "goat/path")

	require.Error(t, err)
	require.Equal(t, `net/http: invalid method "\""`, err.Error())
	require.Nil(t, calls)

	err = client.Clear("\"", "goat/path")

	require.Error(t, err)
	require.Equal(t, `net/http: invalid method "\""`, err.Error())

	client.Port = -1
	err = client.ClearAll()

	require.Error(t, err)
	require.Equal(t, `parse "http://localhost:-1/clear": invalid port ":-1" after host`, err.Error())
}

func TestClientVerifyHttpClientFailure(t *testing.T) {
	client := NewClient()
	client.Close()

	calls, err := client.Verify("GONE", "not/started")

	require.Error(t, err)
	require.Contains(t, err.Error(), `connection refused`)
	require.Nil(t, calls)
}

func TestClientVerifyResponseFailure(t *testing.T) {
	client := NewClient()
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer testServer.Close()
	index := strings.LastIndex(testServer.URL, ":")
	port, err := strconv.ParseInt(testServer.URL[index+1:], 10, 64)
	require.NoError(t, err)
	client.Port = int(port)

	calls, err := client.Verify("GONE", "not/started")

	require.Error(t, err)
	require.Equal(t, `failure to verify calls`, err.Error())
	require.Nil(t, calls)
}

func TestClientVerifyBodyFailure(t *testing.T) {
	client := NewClient()
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode("ydob+dab")
	}))
	defer testServer.Close()
	index := strings.LastIndex(testServer.URL, ":")
	port, err := strconv.ParseInt(testServer.URL[index+1:], 10, 64)
	require.NoError(t, err)
	client.Port = int(port)

	calls, err := client.Verify("BODY", "bad+body")

	require.Error(t, err)
	require.Equal(t, `json: cannot unmarshal string into Go value of type []assured.Call`, err.Error())
	require.Nil(t, calls)
}

func TestClientPathSanitization(t *testing.T) {
	client := NewClient()
	time.Sleep(time.Second)

	require.NoError(t, client.Given(Call{Method: "GET", Path: "///yoyo/path///", StatusCode: http.StatusAccepted}))

	req, err := http.NewRequest(http.MethodGet, client.URL()+"/yoyo/path", nil)
	require.NoError(t, err)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusAccepted, resp.StatusCode)
}

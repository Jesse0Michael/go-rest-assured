package assured

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	kitlog "github.com/go-kit/kit/log"
	"github.com/phayes/freeport"
	"github.com/stretchr/testify/require"
)

func TestClient(t *testing.T) {
	httpClient := &http.Client{}
	logger := kitlog.NewLogfmtLogger(ioutil.Discard)
	port := freeport.GetPort()
	ctx := context.Background()
	client := NewClient(ctx, &port, &logger)

	url := client.URL()

	require.NoError(t, client.Given(call1))
	require.NoError(t, client.Given(call2))
	require.NoError(t, client.Given(call3))

	req, err := http.NewRequest(http.MethodGet, url+"test/assured", bytes.NewReader([]byte(`{"calling":"you"}`)))
	require.NoError(t, err)

	resp, err := httpClient.Do(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	body, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)
	require.Equal(t, []byte(`{"assured": true}`), body)

	req, err = http.NewRequest(http.MethodGet, url+"test/assured", bytes.NewReader([]byte(`{"calling":"again"}`)))
	require.NoError(t, err)

	resp, err = httpClient.Do(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusConflict, resp.StatusCode)
	body, err = ioutil.ReadAll(resp.Body)
	require.NoError(t, err)
	require.Equal(t, []byte("error"), body)

	req, err = http.NewRequest(http.MethodPost, url+"teapot/assured", bytes.NewReader([]byte(`{"calling":"here"}`)))
	require.NoError(t, err)

	resp, err = httpClient.Do(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusTeapot, resp.StatusCode)
	body, err = ioutil.ReadAll(resp.Body)
	require.NoError(t, err)
	require.Equal(t, []byte{}, body)

	calls, err := client.Verify("GET", "test/assured")
	require.NoError(t, err)
	require.Equal(t, []*Call{
		&Call{Method: "GET", Path: "test/assured", StatusCode: 200, Response: []byte(`{"calling":"you"}`)},
		&Call{Method: "GET", Path: "test/assured", StatusCode: 200, Response: []byte(`{"calling":"again"}`)}}, calls)

	calls, err = client.Verify("POST", "teapot/assured")
	require.NoError(t, err)
	require.Equal(t, []*Call{&Call{Method: "POST", Path: "teapot/assured", StatusCode: 200, Response: []byte(`{"calling":"here"}`)}}, calls)

	err = client.Clear("GET", "test/assured")
	require.NoError(t, err)

	calls, err = client.Verify("GET", "test/assured")
	require.NoError(t, err)
	require.Nil(t, calls)

	calls, err = client.Verify("POST", "teapot/assured")
	require.NoError(t, err)
	require.Equal(t, []*Call{&Call{Method: "POST", Path: "teapot/assured", StatusCode: 200, Response: []byte(`{"calling":"here"}`)}}, calls)

	err = client.ClearAll()
	require.NoError(t, err)

	calls, err = client.Verify("GET", "test/assured")
	require.NoError(t, err)
	require.Nil(t, calls)

	calls, err = client.Verify("POST", "teapot/assured")
	require.NoError(t, err)
	require.Nil(t, calls)
}

func TestClientClose(t *testing.T) {
	client := NewDefaultClient()
	client2 := NewDefaultClient()

	require.NoError(t, client.Given(call1))
	require.NoError(t, client2.Given(call1))

	client.Close()
	err := client.Given(call1)

	require.Error(t, err)
	require.Contains(t, err.Error(), `connection refused`)

	client2.Close()
	err = client2.Given(call1)

	require.Error(t, err)
	require.Contains(t, err.Error(), `connection refused`)
}

func TestClientGivenMethodFailure(t *testing.T) {
	client := NewDefaultClient()

	err := client.Given(&Call{Path: "NoMethodMan"})

	require.Error(t, err)
	require.Equal(t, "cannot stub call without Method", err.Error())
}

func TestClientGivenPathFailure(t *testing.T) {
	client := NewDefaultClient()

	err := client.Given(&Call{Method: "GOT"})

	require.Error(t, err)
	require.Equal(t, "cannot stub call without Path", err.Error())
}

func TestClientBadRequestFailure(t *testing.T) {
	client := NewDefaultClient()

	err := client.Given(&Call{Method: "\"", Path: "goat/path"})

	require.Error(t, err)
	require.Equal(t, `net/http: invalid method "\""`, err.Error())

	err = client.Given(&Call{Method: "\"", Path: "goat/path", Response: []byte("goats among men")})

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
	require.Equal(t, `Delete http://localhost:-1/clear: invalid URL port "-1"`, err.Error())
}

func TestClientVerifyHttpClientFailure(t *testing.T) {
	client := NewDefaultClient()
	client.Port = 1

	calls, err := client.Verify("GONE", "not/started")

	require.Error(t, err)
	require.Contains(t, err.Error(), `connection refused`)
	require.Nil(t, calls)
}

func TestClientVerifyResponseFailure(t *testing.T) {
	client := NewDefaultClient()
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
	client := NewDefaultClient()
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode("ydob+dab")
	}))
	defer testServer.Close()
	index := strings.LastIndex(testServer.URL, ":")
	port, err := strconv.ParseInt(testServer.URL[index+1:], 10, 64)
	require.NoError(t, err)
	client.Port = int(port)

	calls, err := client.Verify("BODY", "bad+body")

	require.Error(t, err)
	require.Equal(t, `json: cannot unmarshal string into Go value of type []*assured.Call`, err.Error())
	require.Nil(t, calls)
}

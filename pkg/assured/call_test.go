package assured

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCallKey(t *testing.T) {
	call := Call{
		Path:   "/given/test/assured",
		Method: http.MethodGet,
	}

	require.Equal(t, "GET:/given/test/assured", call.Key())
}

func TestCallKeyNil(t *testing.T) {
	call := Call{}

	require.Equal(t, ":", call.Key())
}

func TestCallString(t *testing.T) {
	call := Call{
		Response: []byte("GO assured is one way to GO"),
	}

	require.Equal(t, "GO assured is one way to GO", call.String())
}

func TestCallStringNil(t *testing.T) {
	call := Call{}

	require.Equal(t, "", call.String())
}

func TestCallUnmarshalNoResponse(t *testing.T) {
	raw := `{
		"path": "teapot/assured", 
		"method": "POST", 
		"status_code": 418, 
		"headers": {"Content-Length": "0", "User-Agent": "Go-http-client/1.1", "Accept-Encoding": "gzip"}
	}`

	call := Call{}
	err := json.Unmarshal([]byte(raw), &call)
	require.NoError(t, err)
	require.Equal(t, *testCall3(), call)
}

func TestCallUnmarshalString(t *testing.T) {
	raw := `{
		"path": "test/assured", 
		"method": "GET", 
		"status_code": 409, 
		"headers": {"Content-Length": "5", "User-Agent": "Go-http-client/1.1", "Accept-Encoding": "gzip"}, 
		"response": "error"
	}`

	call := Call{}
	err := json.Unmarshal([]byte(raw), &call)
	require.NoError(t, err)
	require.Equal(t, *testCall2(), call)
}

func TestCallUnmarshalJSON(t *testing.T) {
	raw := `{
		"path": "test/assured", 
		"method": "GET", 
		"status_code": 200, 
		"headers": {"Content-Length": "17", "User-Agent": "Go-http-client/1.1", "Accept-Encoding": "gzip"}, 
		"query": {"assured": "max"}, 
		"response": "{\"assured\": true}"
	}`

	call := Call{}
	err := json.Unmarshal([]byte(raw), &call)
	require.NoError(t, err)
	require.Equal(t, *testCall1(), call)
}

func TestCallUnmarshalBytes(t *testing.T) {
	raw := `{
		"path": "test/assured", 
		"method": "GET", 
		"status_code": 200, 
		"headers": {"Content-Length": "17", "User-Agent": "Go-http-client/1.1", "Accept-Encoding": "gzip"}, 
		"query": {"assured": "max"}, 
		"response": "eyJhc3N1cmVkIjogdHJ1ZX0="
	}`

	call := Call{}
	err := json.Unmarshal([]byte(raw), &call)
	require.NoError(t, err)
	require.Equal(t, *testCall1(), call)
}

func TestCallUnmarshalFile(t *testing.T) {
	raw := `{
		"path": "test/assured", 
		"method": "GET", 
		"status_code": 200, 
		"headers": {"Content-Length": "17", "User-Agent": "Go-http-client/1.1", "Accept-Encoding": "gzip"}, 
		"query": {"assured": "max"}, 
		"response": "testdata/assured.json"
	}`

	call := Call{}
	err := json.Unmarshal([]byte(raw), &call)
	require.NoError(t, err)
	require.Equal(t, *testCall1(), call)
}

func TestCallUnmarshalCallbacks(t *testing.T) {
	raw := `{
		"path": "test/assured", 
		"method": "GET", 
		"status_code": 200, 
		"headers": {"Content-Length": "17", "User-Agent": "Go-http-client/1.1", "Accept-Encoding": "gzip"},
		"query": {"assured": "max"}, 
		"response": "testdata/assured.json",
		"callbacks": [
			{
				"target": "http://faketarget.com/",
				"method": "POST", 
				"response": "{\"done\": true}", 
				"headers": {"Assured-Callback-Key": "call-key", "Assured-Callback-Target": "http://faketarget.com/"}}
		]
	}`
	expected := *testCall1()
	expected.Callbacks = []Callback{
		{
			Target:   "http://faketarget.com/",
			Response: []byte(`{"done": true}`),
			Method:   "POST",
			Headers:  map[string]string{"Assured-Callback-Key": "call-key", "Assured-Callback-Target": "http://faketarget.com/"},
		},
	}

	call := Call{}
	err := json.Unmarshal([]byte(raw), &call)
	require.NoError(t, err)
	require.Equal(t, expected, call)
}

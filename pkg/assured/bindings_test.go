package assured

import (
	"net/http"
)

// go-rest-assured test vars
// var (
// 	verbs = []string{
// 		http.MethodGet,
// 		http.MethodHead,
// 		http.MethodPost,
// 		http.MethodPut,
// 		http.MethodPatch,
// 		http.MethodDelete,
// 		http.MethodConnect,
// 		http.MethodOptions,
// 	}
// 	fullAssuredCalls = &CallStore{
// 		data: map[string][]*Call{
// 			"GET:test/assured":    {testCall1(), testCall2()},
// 			"POST:teapot/assured": {testCall3()},
// 		},
// 	}
// )

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

// func testCallback() *Call {
// 	return &Call{
// 		Response: []byte(`{"done": true}`),
// 		Method:   "POST",
// 		Headers:  map[string]string{"Assured-Callback-Key": "call-key", "Assured-Callback-Target": "http://faketarget.com/"},
// 	}
// }

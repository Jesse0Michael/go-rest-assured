package assured

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strconv"
	"strings"

	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

const (
	AssuredStatus         = "Assured-Status"
	AssuredCallbackKey    = "Assured-Callback-Key"
	AssuredCallbackTarget = "Assured-Callback-Target"
	AssuredCallbackDelay  = "Assured-Callback-Delay"
)

// StartApplicationHTTPListener creates a Go-routine that has an HTTP listener for the application endpoints
func StartApplicationHTTPListener(root context.Context, errc chan error, settings Settings) {
	go func() {
		ctx, cancel := context.WithCancel(root)
		defer cancel()

		listen, err := net.Listen("tcp", fmt.Sprintf(":%d", settings.Port))
		if err != nil {
			panic(err)
		}

		go func() {
			<-ctx.Done()
			listen.Close()
		}()

		router := createApplicationRouter(ctx, settings)
		settings.Logger.Log("message", fmt.Sprintf("starting go rest assured on port %d", listen.Addr().(*net.TCPAddr).Port))
		errc <- http.Serve(listen, handlers.RecoveryHandler()(router))
	}()
}

// createApplicationRouter sets up the router that will handle all of the application routes
func createApplicationRouter(ctx context.Context, settings Settings) *mux.Router {
	router := mux.NewRouter()
	e := NewAssuredEndpoints(settings)
	assuredMethods := []string{
		http.MethodGet,
		http.MethodHead,
		http.MethodPost,
		http.MethodPut,
		http.MethodPatch,
		http.MethodDelete,
		http.MethodConnect,
		http.MethodOptions,
	}

	router.Handle(
		"/given/{path:.*}",
		kithttp.NewServer(
			e.WrappedEndpoint(e.GivenEndpoint),
			decodeAssuredCall,
			encodeAssuredCall,
			kithttp.ServerErrorLogger(settings.Logger),
			kithttp.ServerAfter(kithttp.SetResponseHeader("Access-Control-Allow-Origin", "*"))),
	).Methods(assuredMethods...)

	router.Handle(
		"/callback",
		kithttp.NewServer(
			e.WrappedEndpoint(e.GivenCallbackEndpoint),
			decodeAssuredCallback,
			encodeAssuredCall,
			kithttp.ServerErrorLogger(settings.Logger),
			kithttp.ServerAfter(kithttp.SetResponseHeader("Access-Control-Allow-Origin", "*"))),
	).Methods(assuredMethods...)

	router.Handle(
		"/when/{path:.*}",
		kithttp.NewServer(
			e.WrappedEndpoint(e.WhenEndpoint),
			decodeAssuredCall,
			encodeAssuredCall,
			kithttp.ServerErrorLogger(settings.Logger),
			kithttp.ServerAfter(kithttp.SetResponseHeader("Access-Control-Allow-Origin", "*"))),
	).Methods(assuredMethods...)

	router.Handle(
		"/verify/{path:.*}",
		kithttp.NewServer(
			e.WrappedEndpoint(e.VerifyEndpoint),
			decodeAssuredCall,
			encodeAssuredCall,
			kithttp.ServerErrorLogger(settings.Logger),
			kithttp.ServerAfter(kithttp.SetResponseHeader("Access-Control-Allow-Origin", "*"))),
	).Methods(assuredMethods...)

	router.Handle(
		"/clear/{path:.*}",
		kithttp.NewServer(
			e.WrappedEndpoint(e.ClearEndpoint),
			decodeAssuredCall,
			encodeAssuredCall,
			kithttp.ServerErrorLogger(settings.Logger),
			kithttp.ServerAfter(kithttp.SetResponseHeader("Access-Control-Allow-Origin", "*"))),
	).Methods(assuredMethods...)

	router.Handle(
		"/clear",
		kithttp.NewServer(
			e.ClearAllEndpoint,
			decodeAssuredCall,
			encodeAssuredCall,
			kithttp.ServerErrorLogger(settings.Logger),
			kithttp.ServerAfter(kithttp.SetResponseHeader("Access-Control-Allow-Origin", "*"))),
	).Methods(http.MethodDelete)

	return router
}

// decodeAssuredCall converts an http request into an assured Call object
func decodeAssuredCall(ctx context.Context, req *http.Request) (interface{}, error) {
	urlParams := mux.Vars(req)
	ac := Call{
		Path:       urlParams["path"],
		Method:     req.Method,
		StatusCode: http.StatusOK,
	}

	// Set status code override
	if statusCode, err := strconv.ParseInt(req.Header.Get(AssuredStatus), 10, 64); err == nil {
		ac.StatusCode = int(statusCode)
	}

	// Set headers
	headers := map[string]string{}
	for key, value := range req.Header {
		headers[key] = value[0]
	}
	ac.Headers = headers

	// Set response body
	if req.Body != nil {
		defer req.Body.Close()
		if bytes, err := ioutil.ReadAll(req.Body); err == nil {
			ac.Response = bytes
		}
	}

	return &ac, nil
}

// decodeAssuredCallback converts an http request into an assured Callback object
func decodeAssuredCallback(ctx context.Context, req *http.Request) (interface{}, error) {
	ac := Call{
		Method:     req.Method,
		StatusCode: http.StatusCreated,
	}

	// Require headers
	if len(req.Header[AssuredCallbackKey]) == 0 {
		return nil, fmt.Errorf("'%s' header required for callback", AssuredCallbackKey)
	}
	if len(req.Header[AssuredCallbackTarget]) == 0 {
		return nil, fmt.Errorf("'%s' header required for callback", AssuredCallbackTarget)
	}

	// Set headers
	headers := map[string]string{}
	for key, value := range req.Header {
		headers[key] = value[0]
	}
	ac.Headers = headers

	// Set response body
	if req.Body != nil {
		defer req.Body.Close()
		if bytes, err := ioutil.ReadAll(req.Body); err == nil {
			ac.Response = bytes
		}
	}

	return &ac, nil
}

// encodeAssuredCall writes the assured Call to the http response as it is intended to be stubbed
func encodeAssuredCall(ctx context.Context, w http.ResponseWriter, i interface{}) error {
	switch resp := i.(type) {
	case *Call:
		w.WriteHeader(resp.StatusCode)
		for key, value := range resp.Headers {
			if !strings.HasPrefix(key, "Assured-") {
				w.Header().Set(key, value)
			}
		}
		w.Write([]byte(resp.String()))
	case []*Call:
		w.Header().Set("Content-Type", "application/json")
		return json.NewEncoder(w).Encode(resp)
	}
	return nil
}

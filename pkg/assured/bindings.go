package assured

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

const (
	AssuredStatus         = "Assured-Status"
	AssuredMethod         = "Assured-Method"
	AssuredDelay          = "Assured-Delay"
	AssuredCallbackKey    = "Assured-Callback-Key"
	AssuredCallbackTarget = "Assured-Callback-Target"
	AssuredCallbackDelay  = "Assured-Callback-Delay"
)

// Serve starts the Rest Assured client to begin listening on the application endpoints
func (c *Client) Serve() error {
	_ = c.logger.Log("message", "starting go rest assured", "port", c.Port)
	if c.tlsCertFile != "" && c.tlsKeyFile != "" {
		return http.ServeTLS(c.listener, handlers.RecoveryHandler()(c.router), c.tlsCertFile, c.tlsKeyFile)
	} else {
		return http.Serve(c.listener, handlers.RecoveryHandler()(c.router))
	}
}

// createApplicationRouter sets up the router that will handle all of the application routes
func (c *Client) createApplicationRouter() *mux.Router {
	router := mux.NewRouter()
	e := NewAssuredEndpoints(c.Options)
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
			kithttp.ServerErrorLogger(c.logger),
			kithttp.ServerAfter(kithttp.SetResponseHeader("Access-Control-Allow-Origin", "*"))),
	).Methods(assuredMethods...)

	router.Handle(
		"/callback",
		kithttp.NewServer(
			e.WrappedEndpoint(e.GivenCallbackEndpoint),
			decodeAssuredCallback,
			encodeAssuredCall,
			kithttp.ServerErrorLogger(c.logger),
			kithttp.ServerAfter(kithttp.SetResponseHeader("Access-Control-Allow-Origin", "*"))),
	).Methods(assuredMethods...)

	router.Handle(
		"/when/{path:.*}",
		kithttp.NewServer(
			e.WrappedEndpoint(e.WhenEndpoint),
			decodeAssuredCall,
			encodeAssuredCall,
			kithttp.ServerErrorLogger(c.logger),
			kithttp.ServerAfter(kithttp.SetResponseHeader("Access-Control-Allow-Origin", "*"))),
	).Methods(assuredMethods...)

	router.Handle(
		"/verify/{path:.*}",
		kithttp.NewServer(
			e.WrappedEndpoint(e.VerifyEndpoint),
			decodeAssuredCall,
			encodeAssuredCall,
			kithttp.ServerErrorLogger(c.logger),
			kithttp.ServerAfter(kithttp.SetResponseHeader("Access-Control-Allow-Origin", "*"))),
	).Methods(assuredMethods...)

	router.Handle(
		"/clear/{path:.*}",
		kithttp.NewServer(
			e.WrappedEndpoint(e.ClearEndpoint),
			decodeAssuredCall,
			encodeAssuredCall,
			kithttp.ServerErrorLogger(c.logger),
			kithttp.ServerAfter(kithttp.SetResponseHeader("Access-Control-Allow-Origin", "*"))),
	).Methods(assuredMethods...)

	router.Handle(
		"/clear",
		kithttp.NewServer(
			e.ClearAllEndpoint,
			decodeAssuredCall,
			encodeAssuredCall,
			kithttp.ServerErrorLogger(c.logger),
			kithttp.ServerAfter(kithttp.SetResponseHeader("Access-Control-Allow-Origin", "*"))),
	).Methods(http.MethodDelete)

	return router
}

// decodeAssuredCall converts an http request into an assured Call object
func decodeAssuredCall(ctx context.Context, req *http.Request) (interface{}, error) {
	urlParams := mux.Vars(req)
	method := req.Method
	if m := req.Header.Get(AssuredMethod); m != "" {
		method = m
	}

	ac := Call{
		Path:       urlParams["path"],
		Method:     method,
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

	// Set query
	query := map[string]string{}
	for key, value := range req.URL.Query() {
		query[key] = value[0]
	}
	ac.Query = query

	// Set response body
	if req.Body != nil {
		defer req.Body.Close()
		if bytes, err := io.ReadAll(req.Body); err == nil {
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
		if bytes, err := io.ReadAll(req.Body); err == nil {
			ac.Response = bytes
		}
	}

	return &ac, nil
}

// encodeAssuredCall writes the assured Call to the http response as it is intended to be stubbed
func encodeAssuredCall(ctx context.Context, w http.ResponseWriter, i interface{}) error {
	switch resp := i.(type) {
	case *Call:
		for key, value := range resp.Headers {
			if !strings.HasPrefix(key, "Assured-") {
				w.Header().Set(key, value)
			}
		}
		w.WriteHeader(resp.StatusCode)
		_, _ = w.Write([]byte(resp.String()))
	case []*Call:
		w.Header().Set("Content-Type", "application/json")
		return json.NewEncoder(w).Encode(resp)
	}
	return nil
}

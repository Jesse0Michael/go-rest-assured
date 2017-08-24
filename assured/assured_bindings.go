package assured

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	kitlog "github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

// StartApplicationHTTPListener creates a Go-routine that has an HTTP listener for the application endpoints
func StartApplicationHTTPListener(port int, logger kitlog.Logger, root context.Context, errc chan error) {
	go func() {
		ctx, cancel := context.WithCancel(root)
		defer cancel()
		router := createApplicationRouter(ctx, logger)
		logger.Log("message", fmt.Sprintf("starting go rest assured on port %d", port))
		errc <- http.ListenAndServe(fmt.Sprintf(":%d", port), handlers.RecoveryHandler()(router))
	}()
}

// createApplicationRouter sets up the router that will handle all of the application routes
func createApplicationRouter(ctx context.Context, logger kitlog.Logger) *mux.Router {
	router := mux.NewRouter()
	e := NewAssuredEndpoints(logger)
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
			kithttp.ServerErrorLogger(logger),
			kithttp.ServerAfter(kithttp.SetResponseHeader("Access-Control-Allow-Origin", "*")),
			kithttp.ServerErrorEncoder(errorEncoder)),
	).Methods(assuredMethods...)

	router.Handle(
		"/when/{path:.*}",
		kithttp.NewServer(
			e.WrappedEndpoint(e.WhenEndpoint),
			decodeAssuredCall,
			encodeAssuredCall,
			kithttp.ServerErrorLogger(logger),
			kithttp.ServerAfter(kithttp.SetResponseHeader("Access-Control-Allow-Origin", "*")),
			kithttp.ServerErrorEncoder(errorEncoder)),
	).Methods(assuredMethods...)

	router.Handle(
		"/verify/{path:.*}",
		kithttp.NewServer(
			e.WrappedEndpoint(e.VerifyEndpoint),
			decodeAssuredCall,
			encodeJSONResponse,
			kithttp.ServerErrorLogger(logger),
			kithttp.ServerAfter(kithttp.SetResponseHeader("Access-Control-Allow-Origin", "*")),
			kithttp.ServerErrorEncoder(errorEncoder)),
	).Methods(assuredMethods...)

	router.Handle(
		"/clear/{path:.*}",
		kithttp.NewServer(
			e.WrappedEndpoint(e.ClearEndpoint),
			decodeAssuredCall,
			encodeAssuredCall,
			kithttp.ServerErrorLogger(logger),
			kithttp.ServerAfter(kithttp.SetResponseHeader("Access-Control-Allow-Origin", "*")),
			kithttp.ServerErrorEncoder(errorEncoder)),
	).Methods(assuredMethods...)

	router.Handle(
		"/clear",
		kithttp.NewServer(
			e.ClearAllEndpoint,
			decodeAssuredCall,
			encodeAssuredCall,
			kithttp.ServerErrorLogger(logger),
			kithttp.ServerAfter(kithttp.SetResponseHeader("Access-Control-Allow-Origin", "*")),
			kithttp.ServerErrorEncoder(errorEncoder)),
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
	if statusCode, err := strconv.ParseInt(req.Header.Get("Assured-Status"), 10, 64); err == nil {
		ac.StatusCode = int(statusCode)
	}

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
	if call, ok := i.(*Call); ok {
		w.WriteHeader(call.StatusCode)
		w.Write([]byte(call.String()))
	}
	return nil
}

// encodeJSONResponse writes to the http response as JSON
func encodeJSONResponse(ctx context.Context, w http.ResponseWriter, i interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(i)
}

func errorEncoder(ctx context.Context, err error, w http.ResponseWriter) {

}

package bindings

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-kit/kit/endpoint"
	kitlog "github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/jesse0michael/go-rest-assured/endpoints"
)

var assuredMethods = []string{
	http.MethodGet,
	http.MethodHead,
	http.MethodPost,
	http.MethodPut,
	http.MethodPatch,
	http.MethodDelete,
	http.MethodConnect,
	http.MethodOptions,
}

// StartApplicationHTTPListener creates a Go-routine that has an HTTP listener for the application endpoints
func StartApplicationHTTPListener(logger kitlog.Logger, root context.Context, errc chan error) {
	go func() {
		ctx, cancel := context.WithCancel(root)
		defer cancel()

		router := createApplicationRouter(ctx, logger)

		errc <- http.ListenAndServe(":9090", handlers.RecoveryHandler()(handlers.CombinedLoggingHandler(kitlog.NewStdlibAdapter(logger), router)))
	}()
}

// createApplicationRouter sets up the router that will handle all of the application routes
func createApplicationRouter(ctx context.Context, logger kitlog.Logger) *mux.Router {
	router := mux.NewRouter()
	e := endpoints.NewAssuredEndpoints(logger)

	router.Handle(
		"/given/{path:.*}",
		kithttp.NewServer(
			buildAssuredEndpoint(e.GivenEndpoint),
			decodeAssuredCall,
			encodeAssuredCall,
			kithttp.ServerErrorLogger(logger),
		)).Methods(assuredMethods...)

	router.Handle(
		"/when/{path:.*}",
		kithttp.NewServer(
			buildAssuredEndpoint(e.WhenEndpoint),
			decodeAssuredCall,
			encodeAssuredCall,
			kithttp.ServerErrorLogger(logger),
		)).Methods(assuredMethods...)

	router.Handle(
		"/then/{path:.*}",
		kithttp.NewServer(
			buildAssuredEndpoint(e.ThenEndpoint),
			decodeAssuredCall,
			encodeAssuredCall,
			kithttp.ServerErrorLogger(logger),
		)).Methods(assuredMethods...)

	router.Handle(
		"/clear/{path:.*}",
		kithttp.NewServer(
			buildAssuredEndpoint(e.ClearEndpoint),
			decodeAssuredCall,
			encodeAssuredCall,
			kithttp.ServerErrorLogger(logger),
		)).Methods(http.MethodDelete)

	router.Handle(
		"/clear",
		kithttp.NewServer(
			e.ClearAllEndpoint,
			decodeAssuredCall,
			encodeAssuredCall,
			kithttp.ServerErrorLogger(logger),
		)).Methods(http.MethodDelete)

	return router
}

func buildAssuredEndpoint(f func(context.Context, *endpoints.AssuredCall) (*endpoints.AssuredCall, error)) endpoint.Endpoint {
	return func(ctx context.Context, i interface{}) (interface{}, error) {
		if call, ok := i.(*endpoints.AssuredCall); ok {
			return f(ctx, call)
		}
		return nil, errors.New("Unable to decode Assured Call")
	}
}

func decodeAssuredCall(ctx context.Context, req *http.Request) (interface{}, error) {
	urlParams := mux.Vars(req)

	ac := endpoints.AssuredCall{
		Path:       urlParams["path"],
		StatusCode: http.StatusOK,
	}
	if statusCode, err := strconv.ParseInt(req.Header.Get("Assured-Status"), 10, 64); err == nil {
		ac.StatusCode = int(statusCode)
	}

	if req.Body != nil {
		defer req.Body.Close()
		json.NewDecoder(req.Body).Decode(ac.Response)
	}

	return ac, nil
}

func encodeAssuredCall(ctx context.Context, w http.ResponseWriter, i interface{}) error {
	if call, ok := i.(*endpoints.AssuredCall); ok {
		w.Header().Set("Content-Type", "application/json") // Content-Type needs to be set before WriteHeader https://golang.org/pkg/net/http/#ResponseWriter
		w.WriteHeader(call.StatusCode)
		return json.NewEncoder(w).Encode(call.Response)
	}
	return nil
}

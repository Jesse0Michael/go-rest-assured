package assured

import (
	"fmt"
	"net"
	"net/http"
)

type Server struct {
	ServerOptions
	listener net.Listener
	router   *http.ServeMux
	calls    *Store[Call]
	records  *Store[Record]
}

// NewServer creates a new go-rest-assured server
func NewServer(opts ...ServerOption) *Server {
	s := Server{
		ServerOptions: DefaultServerOptions,
		calls:         NewStore[Call](),
		records:       NewStore[Record](),
	}
	s.applyOptions(opts...)
	s.router = routes(s.logger, s.calls, s.records, s.httpClient, s.trackRecords)

	var err error
	s.listener, err = net.Listen("tcp", fmt.Sprintf(":%d", s.Port))
	if err != nil {
		s.logger.Error("unable to create http listener", "port", s.Port, "error", err)
	} else {
		s.Port = s.listener.Addr().(*net.TCPAddr).Port
	}

	return &s
}

// Serve starts the Rest Assured client to begin listening on the application endpoints
func (s *Server) Serve() error {
	if s.listener == nil {
		return fmt.Errorf("invalid server")
	}

	if s.tlsCertFile != "" && s.tlsKeyFile != "" {
		return http.ServeTLS(s.listener, s.router, s.tlsCertFile, s.tlsKeyFile)
	} else {
		return http.Serve(s.listener, s.router)
	}
}

// URL returns the url to use to test you stubbed endpoints
func (s *Server) URL() string {
	return s.url()
}

// Close is used to close the running service
func (s *Server) Close() error {
	if s.listener != nil {
		return s.listener.Close()
	}
	return nil
}

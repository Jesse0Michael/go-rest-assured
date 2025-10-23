package assured

import (
	"fmt"
	"log/slog"
	"net/http"
)

var DefaultServerOptions = ServerOptions{
	httpClient:     http.DefaultClient,
	host:           "localhost",
	trackMadeCalls: true,
	logger:         slog.Default(),
}

// ServerOption configures the server behavior.
type ServerOption func(*ServerOptions)

// ServerOptions can be used to configure the rest assured server.
type ServerOptions struct {
	// httpClient used to interact with the rest assured server.
	httpClient *http.Client

	// set the hostname to use in the client. Defaults to localhost.
	host string

	// port for the rest assured server to listen on. Defaults to any available port.
	Port int

	// tlsCertFile is the location of the tls cert for serving https.
	tlsCertFile string

	// tlsKeyFile is the location of the tls key for serving https.
	tlsKeyFile string

	// trackMadeCalls toggles storing the requests made against the rest assured server. Defaults to true.
	trackMadeCalls bool

	// logger to use for logging. Defaults to the default logger.
	logger *slog.Logger
}

func (o *ServerOptions) applyOptions(opts ...ServerOption) {
	for _, opt := range opts {
		opt(o)
	}
}

// WithHTTPClient sets the http client option.
func WithHTTPClient(c http.Client) ServerOption {
	return func(o *ServerOptions) {
		o.httpClient = &c
	}
}

// WithHost sets the host option.
func WithHost(h string) ServerOption {
	return func(o *ServerOptions) {
		if h != "" {
			o.host = h
		}
	}
}

// WithPort sets the port option.
func WithPort(p int) ServerOption {
	return func(o *ServerOptions) {
		if p != 0 {
			o.Port = p
		}
	}
}

// WithTLS sets the tls options.
func WithTLS(cert, key string) ServerOption {
	return func(o *ServerOptions) {
		o.tlsCertFile = cert
		o.tlsKeyFile = key
	}
}

// WithCallTracking sets the trackMadeCalls option.
func WithCallTracking(t bool) ServerOption {
	return func(o *ServerOptions) {
		o.trackMadeCalls = t
	}
}

// WithLogger sets the logger option.
func WithLogger(l *slog.Logger) ServerOption {
	return func(o *ServerOptions) {
		if l != nil {
			o.logger = l
		}
	}
}

// url returns the url to used by the client internally.
func (o *ServerOptions) url() string {
	schema := "http"
	if o.tlsCertFile != "" && o.tlsKeyFile != "" {
		schema = "https"
	}
	return buildURL(schema, o.host, o.Port)
}

func buildURL(schema, host string, port int) string {
	return fmt.Sprintf("%s://%s:%d", schema, host, port)
}

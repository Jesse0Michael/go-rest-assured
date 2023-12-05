package assured

import (
	"log/slog"
	"net/http"
)

var DefaultOptions = Options{
	httpClient:     http.DefaultClient,
	host:           "localhost",
	trackMadeCalls: true,
	logger:         slog.Default(),
}

// Option is a function on that configures rest assured settings
type Option func(*Options)

// Options can be used to configure the rest assured client.
type Options struct {
	// httpClient used to interact with the rest assured server
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

	// logger to use for logging. Defaults the default logger.
	logger *slog.Logger
}

// WithHTTPClient sets the http client option.
func WithHTTPClient(c http.Client) Option {
	return func(o *Options) {
		o.httpClient = &c
	}
}

// WithHost sets the host option.
func WithHost(h string) Option {
	return func(o *Options) {
		if h != "" {
			o.host = h
		}
	}
}

// WithPort sets the port option.
func WithPort(p int) Option {
	return func(o *Options) {
		if p != 0 {
			o.Port = p
		}
	}
}

// WithTLS sets the tls options.
func WithTLS(cert, key string) Option {
	return func(o *Options) {
		o.tlsCertFile = cert
		o.tlsKeyFile = key
	}
}

// WithCallTracking sets the trackMadeCalls option.
func WithCallTracking(t bool) Option {
	return func(o *Options) {
		o.trackMadeCalls = t
	}
}

// WithCallTracking sets the trackMadeCalls option.
func WithLogger(l *slog.Logger) Option {
	return func(o *Options) {
		if l != nil {
			o.logger = l
		}
	}
}

func (o *Options) applyOptions(opts ...Option) {
	for _, opt := range opts {
		opt(o)
	}
}

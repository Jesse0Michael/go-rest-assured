package assured

import "net/http"

var DefaultClientOptions = ClientOptions{
	httpClient: http.DefaultClient,
	baseURL:    "http://localhost",
}

// ClientOption configures the standalone client behavior.
type ClientOption func(*ClientOptions)

// ClientOptions defines just the HTTP transport and base URL for the client.
type ClientOptions struct {
	httpClient *http.Client
	baseURL    string
}

func (o *ClientOptions) applyOptions(opts ...ClientOption) {
	for _, opt := range opts {
		if opt != nil {
			opt(o)
		}
	}
	if o.httpClient == nil {
		o.httpClient = http.DefaultClient
	}
	if o.baseURL == "" {
		o.baseURL = "http://localhost"
	}
}

// WithClientHTTPClient sets the HTTP client used for requests.
func WithClientHTTPClient(c http.Client) ClientOption {
	return func(o *ClientOptions) {
		o.httpClient = &c
	}
}

// WithClientBaseURL sets the base URL for the client.
func WithClientBaseURL(u string) ClientOption {
	return func(o *ClientOptions) {
		if u != "" {
			o.baseURL = u
		}
	}
}

package assured

import "context"

type Assured struct {
	*Client
	*Server
}

// NewAssured creates a new assured instance with both server and client
func NewAssured(opts ...ServerOption) *Assured {
	s := NewServer(opts...)
	clientOpts := []ClientOption{WithClientBaseURL(s.URL())}
	if s.httpClient != nil {
		clientOpts = append(clientOpts, WithClientHTTPClient(*s.httpClient))
	}
	c := NewClient(clientOpts...)
	return &Assured{
		Client: c,
		Server: s,
	}
}

// ServeAssured creates and starts a new assured instance with both server and client
func ServeAssured(ctx context.Context, opts ...ServerOption) (*Assured, error) {
	s := NewServer(opts...)
	clientOpts := []ClientOption{WithClientBaseURL(s.URL())}
	if s.httpClient != nil {
		clientOpts = append(clientOpts, WithClientHTTPClient(*s.httpClient))
	}
	c := NewClient(clientOpts...)

	if err := s.Serve(ctx); err != nil {
		return nil, err
	}

	return &Assured{
		Client: c,
		Server: s,
	}, nil
}

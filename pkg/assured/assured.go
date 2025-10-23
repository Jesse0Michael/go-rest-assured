package assured

type Assured struct {
	*Client
	*Server
}

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

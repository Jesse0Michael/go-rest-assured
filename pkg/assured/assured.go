package assured

type Assured struct {
	*Client
	*Server
}

func NewAssured(opts ...Option) *Assured {
	s := NewServer(opts...)
	c := NewClient(append(opts, WithPort(s.Port))...)
	return &Assured{
		Client: c,
		Server: s,
	}
}

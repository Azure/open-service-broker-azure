package fake

import "context"

// RunFunction describes a function used to provide pluggable runtime behavior
// to the fake implementation of the api.Server interface
type RunFunction func(context.Context) error

// Server is a fake implementation of api.Server used for testing
type Server struct {
	RunBehavior RunFunction
}

// NewServer returns a new, fake implementation of api.Server used for testing
func NewServer() *Server {
	return &Server{
		RunBehavior: defaultRunBehavior,
	}
}

// Start causes the api server to start serving HTTP requests. It will block
// until an error occurs and will return that error.
func (s *Server) Start(ctx context.Context) error {
	return s.RunBehavior(ctx)
}

func defaultRunBehavior(ctx context.Context) error {
	<-ctx.Done()
	return ctx.Err()
}

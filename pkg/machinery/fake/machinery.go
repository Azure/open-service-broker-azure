package fake

import (
	machinery "github.com/RichardKnop/machinery/v1"
	"github.com/RichardKnop/machinery/v1/backends"
	"github.com/RichardKnop/machinery/v1/tasks"
)

// Server is a fake implementation of the machinery.Server interface used for
// testing
type Server struct{}

// Worker is a fake implementation of the machinery.Worker interface used for
// testing
type Worker struct {
	RunBehavior func() error
}

// NewServer returns a fake implementation of machinery.Server used for testing
func NewServer() *Server {
	return &Server{}
}

// RegisterTask registers a task
func (s *Server) RegisterTask(name string, taskFunc interface{}) error {
	return nil
}

// SendTask submits a task for asynchronous completion
func (s *Server) SendTask(*tasks.Signature) (*backends.AsyncResult, error) {
	return nil, nil
}

// NewWorker returns a new *machiner.Worker
func (s *Server) NewWorker(consumerTag string) *machinery.Worker {
	return nil
}

// NewWorker returns a fake implementation of machinery.Worker used for testing
func NewWorker() *Worker {
	return &Worker{
		RunBehavior: defaultRunBehavior,
	}
}

// Launch launches the Worker, deferring to pre-configured behavior
func (w *Worker) Launch() error {
	return w.RunBehavior()
}

func defaultRunBehavior() error {
	for {
	}
}

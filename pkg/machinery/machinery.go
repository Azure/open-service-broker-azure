package machinery

import (
	machinery "github.com/RichardKnop/machinery/v1"
	"github.com/RichardKnop/machinery/v1/backends"
	"github.com/RichardKnop/machinery/v1/tasks"
)

// Server is an interface created ex-post-facto that is implemented by
// the existing *github.com/RichardKnop/machinery/v1.Server type. By having
// broker components code to this interface, it is easier to swap in an
// alternate implementation to facilitate testing.
type Server interface {
	RegisterTask(name string, taskFunc interface{}) error
	SendTask(*tasks.Signature) (*backends.AsyncResult, error)
	NewWorker(consumerTag string) *machinery.Worker
}

// Worker is an interface created ex-post-facto that is implemented by
// the existing *github.com/RichardKnop/machinery/v1.Worker type. By having
// broker components code to this interface, it is easier to swap in an
// alternate implementation to facilitate testing.
type Worker interface {
	Launch() error
}

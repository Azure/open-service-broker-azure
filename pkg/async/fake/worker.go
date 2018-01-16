package fake

import (
	"context"

	"github.com/Azure/open-service-broker-azure/pkg/async/model"
)

// Worker is a fake implementation of async.Worker used for testing
type Worker struct {
	RunBehavior RunFn
}

// NewWorker returns a new, fake implementation of async.Worker used for testing
func NewWorker() *Worker {
	return &Worker{
		RunBehavior: defaultWorkerRunBehavior,
	}
}

// GetID returns the worker's ID
func (w *Worker) GetID() string {
	return "fake-worker"
}

// RegisterJob registers a new Job with the worker
func (w *Worker) RegisterJob(name string, fn model.JobFn) error {
	return nil
}

// Run causes the worker to process tasks
func (w *Worker) Run(ctx context.Context) error {
	return w.RunBehavior(ctx)
}

func defaultWorkerRunBehavior(ctx context.Context) error {
	<-ctx.Done()
	return ctx.Err()
}

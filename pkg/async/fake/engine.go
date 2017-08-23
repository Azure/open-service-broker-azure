package fake

import (
	"context"

	"github.com/Azure/azure-service-broker/pkg/async/model"
)

// Engine is a fake implementation of async.Engine used for testing
type Engine struct {
	SubmittedTasks map[string]model.Task
	RunBehavior    RunFunction
}

// NewEngine returns a new, fake implementation of async.Engine used for testing
func NewEngine() *Engine {
	return &Engine{
		SubmittedTasks: make(map[string]model.Task),
		RunBehavior:    defaultEngineRunBehavior,
	}
}

// RegisterJob registers a new Job with the async engine
func (e *Engine) RegisterJob(name string, fn model.JobFunction) error {
	return nil
}

// SubmitTask submits an idempotent task to the async engine for reliable,
// asynchronous completion
func (e *Engine) SubmitTask(task model.Task) error {
	e.SubmittedTasks[task.GetID()] = task
	return nil
}

// Start causes the async engine to begin executing queued tasks
func (e *Engine) Start(ctx context.Context) error {
	return e.RunBehavior(ctx)
}

func defaultEngineRunBehavior(ctx context.Context) error {
	<-ctx.Done()
	return ctx.Err()
}

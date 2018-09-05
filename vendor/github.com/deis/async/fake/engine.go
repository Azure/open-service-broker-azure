package fake

import (
	"context"

	"github.com/deis/async"
)

// RunFn describes a function used to provide pluggable runtime behavior to
// the fake async engine
type RunFn func(context.Context) error

// Engine is a fake implementation of async.Engine used for testing
type Engine struct {
	SubmittedTasks map[string]async.Task
	RunBehavior    RunFn
}

// NewEngine returns a new, fake implementation of async.Engine used for testing
func NewEngine() *Engine {
	return &Engine{
		SubmittedTasks: make(map[string]async.Task),
		RunBehavior:    defaultEngineRunBehavior,
	}
}

// RegisterJob registers a new Job with the async engine
func (e *Engine) RegisterJob(name string, fn async.JobFn) error {
	return nil
}

// SubmitTask submits an idempotent task to the async engine for reliable,
// asynchronous completion
func (e *Engine) SubmitTask(task async.Task) error {
	e.SubmittedTasks[task.GetID()] = task
	return nil
}

// Run causes the async engine to carry out all of its functions. It blocks
// until a fatal error is encountered or the context passed to it has been
// canceled. Run always returns a non-nil error.
func (e *Engine) Run(ctx context.Context) error {
	return e.RunBehavior(ctx)
}

func defaultEngineRunBehavior(ctx context.Context) error {
	<-ctx.Done()
	return ctx.Err()
}

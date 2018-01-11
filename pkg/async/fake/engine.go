package fake

import (
	"context"

	"github.com/Azure/open-service-broker-azure/pkg/async/model"
)

// Engine is a fake implementation of async.Engine used for testing
type Engine struct {
	SubmittedTasks map[string]model.Task
	DelayedTasks   map[string][]model.Task
	RunBehavior    RunFunction
}

// NewEngine returns a new, fake implementation of async.Engine used for testing
func NewEngine() *Engine {
	return &Engine{
		SubmittedTasks: make(map[string]model.Task),
		DelayedTasks:   make(map[string][]model.Task),
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

// SubmitDelayedTask submits an idempotent task to the async engine for
// reliable, asynchronous completion, in a delayed state. The task will be
// automatically started on a periodic basis in the future or can be
// started by a client. The task will be stored in a queue named by the
// identifier parameter
func (e *Engine) SubmitDelayedTask(identifier string, task model.Task) error {
	tasks := e.DelayedTasks[identifier]
	e.DelayedTasks[identifier] = append(tasks, task)
	return nil
}

// StartDelayedTasks will transfer all delayed tasks from the queue
// identified by the identifier parameter to the main worker queue for
// processing. The resumer will also start these on a periodic basis if
// they have not been triggered by a client
func (e *Engine) StartDelayedTasks(identifier string) error {
	tasks := e.DelayedTasks[identifier]
	for _, task := range tasks {
		e.SubmittedTasks[task.GetID()] = task
	}
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

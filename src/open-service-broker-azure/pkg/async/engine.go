package async

import (
	"context"
)

// Engine is an interface for a broker-specifc framework for submitting and
// asynchronously completing provisioning and deprovisioning tasks.
type Engine interface {
	// RegisterJob registers a new Job with the async engine
	RegisterJob(name string, fn JobFn) error
	// SubmitTask submits an idempotent task to the async engine for reliable,
	// asynchronous completion
	SubmitTask(Task) error
	// Run causes the async engine to carry out all of its functions. It blocks
	// until a fatal error is encountered or the context passed to it has been
	// canceled. Run always returns a non-nil error.
	Run(context.Context) error
}

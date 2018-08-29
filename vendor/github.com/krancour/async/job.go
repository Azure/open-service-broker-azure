package async

import "context"

// JobFn is the signature for functions that workers can call to asynchronously
// execute a job
type JobFn func(ctx context.Context, task Task) ([]Task, error)

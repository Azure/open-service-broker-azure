package model

import "context"

// JobFn is the signature for functions that workers can call to asynchronously
// execute a job
type JobFn func(ctx context.Context, args map[string]string) error

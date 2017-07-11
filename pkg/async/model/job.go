package model

import "context"

// JobFunction is the signature for functions that workers can call to
// asynchronously execute a job
type JobFunction func(ctx context.Context, args map[string]string) error

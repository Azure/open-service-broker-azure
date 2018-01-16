package fake

import "context"

// Cleaner is a fake implementation of async.Cleaner used for testing
type Cleaner struct {
	RunBehavior RunFn
}

// NewCleaner returns a new, fake implementation of async.Cleaner used for
// testing
func NewCleaner() *Cleaner {
	return &Cleaner{
		RunBehavior: defaultCleanerRunBehavior,
	}
}

// Run causes the cleaner to clean up after dead workers
func (c *Cleaner) Run(ctx context.Context) error {
	return c.RunBehavior(ctx)
}

func defaultCleanerRunBehavior(ctx context.Context) error {
	<-ctx.Done()
	return ctx.Err()
}

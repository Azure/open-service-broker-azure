package fake

import "context"

// Cleaner is a fake implementation of async.Cleaner used for testing
type Cleaner struct {
	RunBehavior RunFunction
}

// NewCleaner returns a new, fake implementation of async.Cleaner used for
// testing
func NewCleaner() *Cleaner {
	return &Cleaner{
		RunBehavior: defaultCleanerRunBehavior,
	}
}

// Clean causes the cleaner to begin cleaning up after dead workers
func (c *Cleaner) Clean(ctx context.Context) error {
	return c.RunBehavior(ctx)
}

func defaultCleanerRunBehavior(ctx context.Context) error {
	<-ctx.Done()
	return ctx.Err()
}

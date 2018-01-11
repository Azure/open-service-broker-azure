package fake

import "context"

// Resumer is a fake implementation of async.Resumer used for testing
type Resumer struct {
	RunBehavior RunFunction
}

// NewResumer returns a new, fake implementation of async.Resumer used for
// testing
func NewResumer() *Resumer {
	return &Resumer{
		RunBehavior: defaultResumerRunBehavior,
	}
}

// Resume moves delayed tasks from watched queues to
// the main worker queue
func (r *Resumer) Resume(ctx context.Context) error {
	return r.RunBehavior(ctx)
}

// Watch adds a queue to the set of watched queues
func (r *Resumer) Watch(q string) {
}

func defaultResumerRunBehavior(ctx context.Context) error {
	<-ctx.Done()
	return ctx.Err()
}

package fake

import "context"

// Heart is a fake implementation of async.Heart used for testing
type Heart struct {
	RunBehavior RunFn
}

// NewHeart returns a new, fake implementation of async.Heart used for testing
func NewHeart() *Heart {
	return &Heart{
		RunBehavior: defaultHeartRunBehavior,
	}
}

// Beat sends a single heartbeat
func (*Heart) Beat() error {
	return nil
}

// Run sends heartbeats at regular intervals
func (h *Heart) Run(ctx context.Context) error {
	return h.RunBehavior(ctx)
}

func defaultHeartRunBehavior(ctx context.Context) error {
	<-ctx.Done()
	return ctx.Err()
}

package redis

import (
	"context"
	"testing"
	"time"

	"github.com/Azure/open-service-broker-azure/pkg/async/redis/fake"
	"github.com/stretchr/testify/assert"
)

func TestNewEnginesHaveUniqueWorkerIDs(t *testing.T) {
	// Create two workers
	e1 := NewEngine(redisClient).(*engine)
	e2 := NewEngine(redisClient).(*engine)

	// Assert that their workerIDs are at least different from one another
	assert.NotEqual(t, e1.workerID, e2.workerID)
}

func TestEngineRunBlocksUntilCleanerStops(t *testing.T) {
	e := NewEngine(redisClient).(*engine)

	// Override the engine's default clean function so it just returns an error
	e.clean = func(context.Context, string, string, string) error {
		return errSome
	}

	// Override the engine's default runHeart function so it just communicates
	// when the context it was passed has been canceled
	contextCanceledCh := make(chan struct{})
	e.runHeart = func(ctx context.Context) error {
		<-ctx.Done()
		close(contextCanceledCh)
		return ctx.Err()
	}

	// Create a fake worker that will just run until the context it was passed is
	// canceled
	w := fake.NewWorker()
	w.RunBehavior = func(ctx context.Context) error {
		<-ctx.Done()
		return ctx.Err()
	}

	// Make the engine use the fake worker
	e.worker = w

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Call Run in a goroutine. If it never unblocks, as we hope it does, we don't
	// want the test to stall.
	errCh := make(chan error)
	go func() {
		errCh <- e.Run(ctx)
	}()

	// Assert that the error returned from the Run function wraps the error that
	// the fake heart generated
	select {
	case err := <-errCh:
		assert.Equal(t, &errCleanerStopped{err: errSome}, err)
	case <-time.After(time.Second):
		assert.Fail(t, "an error should have been received, but wasn't")
	}

	// Assert that the context got canceled. It's helpful to know that when the
	// cleaner stops, the rest of the engine components are also signaled to shut
	// down.
	select {
	case <-contextCanceledCh:
	case <-time.After(time.Second):
		assert.Fail(t, "context should have been canceled, but it was not")
	}
}

func TestEngineRunBlocksUntilHeartStops(t *testing.T) {
	e := NewEngine(redisClient).(*engine)

	// Override the worker's default clean function so it just communicates when
	// the context it was passed has been canceled
	contextCanceledCh := make(chan struct{})
	e.clean = func(ctx context.Context, _ string, _ string, _ string) error {
		<-ctx.Done()
		close(contextCanceledCh)
		return ctx.Err()
	}

	// Override the engine's default runHeart function so it just returns an error
	e.runHeart = func(context.Context) error {
		return errSome
	}

	// Create a fake worker that will just run until the context it was passed is
	// canceled
	w := fake.NewWorker()
	w.RunBehavior = func(ctx context.Context) error {
		<-ctx.Done()
		return ctx.Err()
	}

	// Make the engine use the fake worker
	e.worker = w

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Call Run in a goroutine. If it never unblocks, as we hope it does, we don't
	// want the test to stall.
	errCh := make(chan error)
	go func() {
		errCh <- e.Run(ctx)
	}()

	// Assert that the error returned from the Run function wraps the error that
	// the fake heart generated
	select {
	case err := <-errCh:
		assert.Equal(t, &errHeartStopped{workerID: e.workerID, err: errSome}, err)
	case <-time.After(time.Second):
		assert.Fail(t, "an error should have been received, but wasn't")
	}

	// Assert that the context got canceled. It's helpful to know that when the
	// cleaner stops, the rest of the engine components are also signaled to shut
	// down.
	select {
	case <-contextCanceledCh:
	case <-time.After(time.Second):
		assert.Fail(t, "context should have been canceled, but it was not")
	}
}

func TestEngineRunBlocksUntilWorkerStops(t *testing.T) {
	e := NewEngine(redisClient).(*engine)

	// Override the worker's default clean function so it just communicates when
	// the context it was passed has been canceled
	contextCanceledCh := make(chan struct{})
	e.clean = func(ctx context.Context, _ string, _ string, _ string) error {
		<-ctx.Done()
		close(contextCanceledCh)
		return ctx.Err()
	}

	// Override the engine's default runHeart function so it just blocks until the
	// context it was passed has been canceled
	e.runHeart = func(ctx context.Context) error {
		<-ctx.Done()
		return ctx.Err()
	}

	// Create a fake worker that will just return an error when it runs
	w := fake.NewWorker()
	w.RunBehavior = func(context.Context) error {
		return errSome
	}

	// Make the engine use the fake worker
	e.worker = w

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Call Run in a goroutine. If it never unblocks, as we hope it does, we don't
	// want the test to stall.
	errCh := make(chan error)
	go func() {
		errCh <- e.Run(ctx)
	}()

	// Assert that the error returned from the Run function wraps the error that
	// the fake heart generated
	select {
	case err := <-errCh:
		assert.Equal(t, &errWorkerStopped{workerID: w.GetID(), err: errSome}, err)
	case <-time.After(time.Second):
		assert.Fail(t, "an error should have been received, but wasn't")
	}

	// Assert that the context got canceled. It's helpful to know that when the
	// worker stops, the rest of the engine components are also signaled to shut
	// down.
	select {
	case <-contextCanceledCh:
	case <-time.After(time.Second):
		assert.Fail(t, "context should have been canceled, but it was not")
	}
}

func TestEngineRunRespondsToContextCanceled(t *testing.T) {
	e := NewEngine(redisClient).(*engine)

	// Override the worker's default clean function so it just communicates when
	// the context it was passed has been canceled
	contextCanceledCh := make(chan struct{})
	e.clean = func(ctx context.Context, _ string, _ string, _ string) error {
		<-ctx.Done()
		close(contextCanceledCh)
		return ctx.Err()
	}

	// Override the engine's default runHeart function so it just blocks until the
	// context it was passed has been canceled
	e.runHeart = func(ctx context.Context) error {
		<-ctx.Done()
		return ctx.Err()
	}

	// Create a fake worker that will just run until the context it was passed
	// is canceled
	w := fake.NewWorker()
	w.RunBehavior = func(ctx context.Context) error {
		<-ctx.Done()
		return ctx.Err()
	}

	// Make the engine use the fake worker
	e.worker = w

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Call Run in a goroutine. If it never unblocks, as we hope it does, we don't
	// want the test to stall.
	errCh := make(chan error)
	go func() {
		errCh <- e.Run(ctx)
	}()

	cancel()

	// Assert that the error returned from Run indicates that the context was
	// canceled
	select {
	case err := <-errCh:
		assert.Equal(t, ctx.Err(), err)
	case <-time.After(time.Second):
		assert.Fail(
			t,
			"a context canceled error should have been returned, but wasn't",
		)
	}
}

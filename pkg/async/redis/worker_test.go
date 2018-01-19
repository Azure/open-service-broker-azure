package redis

import (
	"context"
	"testing"
	"time"

	"github.com/Azure/open-service-broker-azure/pkg/async/redis/fake"
	"github.com/stretchr/testify/assert"
)

func TestNewWorkersHaveUniqueIDs(t *testing.T) {
	// Create two workers
	w1 := newWorker(redisClient).(*worker)
	w2 := newWorker(redisClient).(*worker)

	// Assert that their IDs are at least different from one another
	assert.NotEqual(t, w1.id, w2.id)
}

// TestWorkerRunBlocksUntilHeartStops tests what happens when a worker's heart
// stops beating.
func TestWorkerRunBlocksUntilHeartStops(t *testing.T) {
	// Use a fake heart
	h := fake.NewHeart()

	// Specify the fake heart's runtime behavior should just return an error
	h.RunBehavior = func(context.Context) error {
		return errSome
	}

	w := newWorker(redisClient).(*worker)

	// Make the worker use the fake heart
	w.heart = h

	// Override the worker's default receivePendingTasks function so it just
	// communicates when the context it was passed has been canceled
	contextCanceledCh := make(chan struct{})
	w.receivePendingTasks = func(
		ctx context.Context,
		_ string,
		_ string,
		_ chan []byte,
		_ chan error,
	) {
		<-ctx.Done()
		close(contextCanceledCh)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Call Run in a goroutine. If it never unblocks, as we hope it does, we don't
	// want the test to stall.
	errCh := make(chan error)
	go func() {
		errCh <- w.Run(ctx)
	}()

	// Assert that the error returned from the Run function wraps the error that
	// the fake heart generated
	select {
	case err := <-errCh:
		assert.Equal(t, &errHeartStopped{workerID: w.id, err: errSome}, err)
	case <-time.After(time.Second):
		assert.Fail(t, "an error should have been received, but wasn't")
	}

	// Assert that the context got canceled. It's helpful to know that when the
	// heart stops, the rest of the worker components are also signaled to shut
	// down.
	select {
	case <-contextCanceledCh:
	case <-time.After(time.Second):
		assert.Fail(t, "context should have been canceled, but it was not")
	}
}

// TestWorkerRunBlocksUntilPendingReceiverStops tests what happens when a
// worker's goroutine that receives pending tasks stops running.
func TestWorkerRunBlocksUntilPendingReceiverStops(t *testing.T) {
	// Use a fake heart
	h := fake.NewHeart()

	// Specify the fake heart's runtime behavior should just communicates when the
	// context it was passed has been canceled
	contextCanceledCh := make(chan struct{})
	h.RunBehavior = func(ctx context.Context) error {
		<-ctx.Done()
		close(contextCanceledCh)
		return ctx.Err()
	}

	w := newWorker(redisClient).(*worker)

	// Make the worker use the fake heart
	w.heart = h

	// Override the worker's default receivePendingTasks function so it just
	// returns an error and then blocks until the context it was passed is
	// canceled
	w.receivePendingTasks = func(
		ctx context.Context,
		_ string,
		_ string,
		_ chan []byte,
		errCh chan error,
	) {
		select {
		case errCh <- errSome:
		case <-ctx.Done():
		}
		<-ctx.Done()
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Call Run in a goroutine. If it never unblocks, as we hope it does, we don't
	// want the test to stall.
	errCh := make(chan error)
	go func() {
		errCh <- w.Run(ctx)
	}()

	// Assert that the error received from the Run function wraps the error that
	// the overridden receivePendingTasks function generated
	select {
	case err := <-errCh:
		assert.Equal(
			t,
			&errReceiverStopped{
				workerID:  w.id,
				queueName: pendingTaskQueueName,
				err:       errSome,
			},
			err,
		)
	case <-time.After(time.Second):
		assert.Fail(t, "an error should have been received, but wasn't")
	}

	// Assert that the context got canceled. It's helpful to know that when the
	// pending task receiver stops, the rest of the worker components are also
	// signaled to shut down.
	select {
	case <-contextCanceledCh:
	case <-time.After(time.Second):
		assert.Fail(t, "context should have been canceled, but it was not")
	}
}

// TestWorkerRunBlocksUntilExecuteTasksStops tests what happens when a worker's
// goroutine that executes pending tasks stops.
func TestWorkerRunBlocksUntilExecuteTasksStop(t *testing.T) {
	// Use a fake heart
	h := fake.NewHeart()

	// Specify the fake heart's runtime behavior should just communicates when the
	// context it was passed has been canceled
	contextCanceledCh := make(chan struct{})
	h.RunBehavior = func(ctx context.Context) error {
		<-ctx.Done()
		close(contextCanceledCh)
		return ctx.Err()
	}

	w := newWorker(redisClient).(*worker)

	// Make the worker use the fake heart
	w.heart = h

	// Override the worker's default executeTasks function so it just sends an
	// error and then blocks until the context it was passed is canceled
	w.executeTasks = func(
		ctx context.Context,
		_ chan []byte,
		_ string,
		_ string,
		errCh chan error,
	) {
		select {
		case errCh <- errSome:
		case <-ctx.Done():
		}
		<-ctx.Done()
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Call Run in a goroutine. If it never unblocks, as we hope it does, we don't
	// want the test to stall.
	errCh := make(chan error)
	go func() {
		errCh <- w.Run(ctx)
	}()

	// Assert that the error received from the Run function wraps the error that
	// the overridden executeTasks function generated
	select {
	case err := <-errCh:
		assert.Equal(t, &errTaskExecutorStopped{workerID: w.id, err: errSome}, err)
	case <-time.After(time.Second):
		assert.Fail(t, "an error should have been received, but wasn't")
	}

	// Assert that the context got canceled. It's helpful to know that when the
	// task executor stops, the rest of the worker components are also signaled to
	// shut down.
	select {
	case <-contextCanceledCh:
	case <-time.After(time.Second):
		assert.Fail(t, "context should have been canceled, but it was not")
	}
}

// TestWorkerRunBlocksUntilDeferredReceiverStops tests what happens when a
// worker's goroutine that receives deferred tasks stops.
func TestWorkerRunBlocksUntilDeferredReceiverStops(t *testing.T) {
	// Use a fake heart
	h := fake.NewHeart()

	// Specify the fake heart's runtime behavior should just communicates when the
	// context it was passed has been canceled
	contextCanceledCh := make(chan struct{})
	h.RunBehavior = func(ctx context.Context) error {
		<-ctx.Done()
		close(contextCanceledCh)
		return ctx.Err()
	}

	w := newWorker(redisClient).(*worker)

	// Make the worker use the fake heart
	w.heart = h

	// Override the worker's default receiveDeferredTasks function so it just
	// sends an error and then blocks until the context it was passed is canceled
	w.receiveDeferredTasks = func(
		ctx context.Context,
		_ string,
		_ string,
		_ chan []byte,
		errCh chan error,
	) {
		select {
		case errCh <- errSome:
		case <-ctx.Done():
		}
		<-ctx.Done()
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Call Run in a goroutine. If it never unblocks, as we hope it does, we don't
	// want the test to stall.
	errCh := make(chan error)
	go func() {
		errCh <- w.Run(ctx)
	}()

	// Assert that the error received from the Run function wraps the error that
	// the overridden receiveDeferredTasks function generated
	select {
	case err := <-errCh:
		assert.Equal(
			t,
			&errReceiverStopped{
				workerID:  w.id,
				queueName: deferredTaskQueueName,
				err:       errSome,
			},
			err,
		)
	case <-time.After(time.Second):
		assert.Fail(t, "an error should have been received, but wasn't")
	}

	// Assert that the context got canceled. It's helpful to know that when the
	// deferred task receiver stops, the rest of the worker components are also
	// signaled to shut down.
	select {
	case <-contextCanceledCh:
	case <-time.After(time.Second):
		assert.Fail(t, "context should have been canceled, but it was not")
	}
}

// TestWorkerRunBlocksUntilWatchDeferredTaskErrors tests what happens when a
// worker's goroutine that watches a deferred task errors.
func TestWorkerRunBlocksUntilWatchDeferredTaskErrors(t *testing.T) {
	// Use a fake heart
	h := fake.NewHeart()

	// Specify the fake heart's runtime behavior should just communicates when the
	// context it was passed has been canceled
	contextCanceledCh := make(chan struct{})
	h.RunBehavior = func(ctx context.Context) error {
		<-ctx.Done()
		close(contextCanceledCh)
		return ctx.Err()
	}

	w := newWorker(redisClient).(*worker)

	// Make the worker use the fake heart
	w.heart = h

	// Override the worker's default receiveDeferredTasks function so it just
	// sends a result (to trigger a new goroutine running the watchDeferredTask
	// function) and then it blocks until the context it was passed is canceled.
	w.receiveDeferredTasks = func(
		ctx context.Context,
		_ string,
		_ string,
		retCh chan []byte,
		_ chan error,
	) {
		select {
		case retCh <- []byte{}: // A dummy value is fine
		case <-ctx.Done():
		}
		<-ctx.Done()
	}

	// Override the worker's default watchDeferredTask function so it just sends
	// an error and then blocks until the context it was passed is canceled
	w.watchDeferredTask = func(
		ctx context.Context,
		_ []byte,
		_ string,
		errCh chan error,
	) {
		select {
		case errCh <- errSome:
		case <-ctx.Done():
		}
		<-ctx.Done()
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Call Run in a goroutine. If it never unblocks, as we hope it does, we don't
	// want the test to stall.
	errCh := make(chan error)
	go func() {
		errCh <- w.Run(ctx)
	}()

	// Assert that the error received from the Run function wraps the error that
	// the overridden watchDeferredTask function generated
	select {
	case err := <-errCh:
		assert.Equal(
			t,
			&errDeferredTaskWatcherStopped{workerID: w.id, err: errSome},
			err,
		)
	case <-time.After(time.Second):
		assert.Fail(t, "an error should have been received, but wasn't")
	}

	// Assert that the context got canceled. It's helpful to know that when the
	// a deferred task watcher errors, the rest of the worker components are also
	// signaled to shut down.
	select {
	case <-contextCanceledCh:
	case <-time.After(time.Second):
		assert.Fail(t, "context should have been canceled, but it was not")
	}
}

// TestWorkerRunRespondsToContextCanceled tests that canceling the context
// passed to the Run function causes the Run function to return.
func TestWorkerRunRespondsToContextCanceled(t *testing.T) {
	// Use a fake heart
	h := fake.NewHeart()

	// Specify the fake heart's runtime behavior should just block until its
	// context is canceled
	h.RunBehavior = func(ctx context.Context) error {
		<-ctx.Done()
		return ctx.Err()
	}

	w := newWorker(redisClient).(*worker)

	// Make the worker use the fake heart
	w.heart = h

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Call Run in a goroutine. If it never unblocks, as we hope it does, we don't
	// want the test to stall.
	errCh := make(chan error)
	go func() {
		errCh <- w.Run(ctx)
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

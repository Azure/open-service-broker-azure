package redis

import (
	"context"
	"testing"
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

// TestWorkerRunBlocksUntilPendingReceiverStops tests what happens when a
// worker's goroutine that receives pending tasks stops running.
func TestWorkerRunBlocksUntilPendingReceiverStops(t *testing.T) {
	w := newWorker(redisClient, uuid.NewV4().String()).(*worker)

	// Override the worker's default receivePendingTasks function so it just
	// returns an error and then blocks until the context it was passed is
	// canceled. It will also communicate when the context it was passed is
	// canceled.
	contextCanceledCh := make(chan struct{})
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
	w := newWorker(redisClient, uuid.NewV4().String()).(*worker)

	// Override the worker's default executeTasks function so it just sends an
	// error and then blocks until the context it was passed is canceled. It will
	// also communicate when the context it was passed is canceled.
	contextCanceledCh := make(chan struct{})
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
		contextCanceledCh <- struct{}{}
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
	w := newWorker(redisClient, uuid.NewV4().String()).(*worker)

	// Override the worker's default receiveDeferredTasks function so it just
	// sends an error and then blocks until the context it was passed is canceled.
	// It will also communicate when the context it was passed is canceled.
	contextCanceledCh := make(chan struct{})
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
	w := newWorker(redisClient, uuid.NewV4().String()).(*worker)

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
	// an error and then blocks until the context it was passed is canceled.
	// It will also communicate when the context it was passed is canceled.
	contextCanceledCh := make(chan struct{})
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
	w := newWorker(redisClient, uuid.NewV4().String()).(*worker)

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

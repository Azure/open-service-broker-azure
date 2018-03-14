package redis

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewEnginesHaveUniqueWorkerIDs(t *testing.T) {
	// Create two engines
	e1 := getTestEngine()
	e2 := getTestEngine()

	// Assert that their workerIDs are at least different from one another
	assert.NotEqual(t, e1.workerID, e2.workerID)
}

func TestRunBlocksUntilCleanReturnsError(t *testing.T) {
	e := getTestEngine()

	// Override the engine's clean function so it just returns an error
	e.clean = func(context.Context, string, string, string, time.Duration) error {
		return errSome
	}

	// Override the engine's runHeart function so it just communicates when the
	// context it was passed has been canceled
	contextCanceledCh := make(chan struct{})
	e.runHeart = func(ctx context.Context, _ time.Duration) error {
		<-ctx.Done()
		close(contextCanceledCh)
		return ctx.Err()
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Call Run in a goroutine. If it never unblocks, as we hope it does, we don't
	// want the test to stall.
	errCh := make(chan error)
	go func() {
		errCh <- e.Run(ctx)
	}()

	// Assert that the error returned from the Run function wraps the error that
	// the overridden clean function returned
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

func TestRunBlocksUntilRunHeartReturnsError(t *testing.T) {
	e := getTestEngine()

	// Override the engine's clean function so it just communicates when the
	// context it was passed has been canceled
	contextCanceledCh := make(chan struct{})
	e.clean = func(
		ctx context.Context,
		_ string,
		_ string,
		_ string,
		_ time.Duration,
	) error {
		<-ctx.Done()
		close(contextCanceledCh)
		return ctx.Err()
	}

	// Override the engine's runHeart function so it just returns an error
	e.runHeart = func(context.Context, time.Duration) error {
		return errSome
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Call Run in a goroutine. If it never unblocks, as we hope it does, we don't
	// want the test to stall.
	errCh := make(chan error)
	go func() {
		errCh <- e.Run(ctx)
	}()

	// Assert that the error returned from the Run function wraps the error that
	// the overridden runHeart function returned
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

func TestRunBlocksUntilReceivePendingTasksSendsError(t *testing.T) {
	e := getTestEngine()

	// Override the engine's clean function so it just communicates when the
	// context it was passed has been canceled
	contextCanceledCh := make(chan struct{})
	e.clean = func(
		ctx context.Context,
		_ string,
		_ string,
		_ string,
		_ time.Duration,
	) error {
		<-ctx.Done()
		close(contextCanceledCh)
		return ctx.Err()
	}

	// Override the engine's receivePendingTasks function so it just sends an
	// error
	e.receivePendingTasks = func(
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
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Call Run in a goroutine. If it never unblocks, as we hope it does, we don't
	// want the test to stall.
	errCh := make(chan error)
	go func() {
		errCh <- e.Run(ctx)
	}()

	// Assert that the error returned from the Run function wraps the error that
	// the overridden receivePendingTasks function sent
	select {
	case err := <-errCh:
		assert.Equal(
			t,
			&errReceiverStopped{
				workerID:  e.workerID,
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

func TestRunBlocksUntilExecuteTasksSendsError(t *testing.T) {
	e := getTestEngine()

	// Override the engine's clean function so it just communicates when the
	// context it was passed has been canceled
	contextCanceledCh := make(chan struct{})
	e.clean = func(
		ctx context.Context,
		_ string,
		_ string,
		_ string,
		_ time.Duration,
	) error {
		<-ctx.Done()
		close(contextCanceledCh)
		return ctx.Err()
	}

	// Override the engine's executeTasks function so it just sends an error
	e.executeTasks = func(
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
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Call Run in a goroutine. If it never unblocks, as we hope it does, we don't
	// want the test to stall.
	errCh := make(chan error)
	go func() {
		errCh <- e.Run(ctx)
	}()

	// Assert that the error returned from the Run function wraps the error that
	// the overridden executeTasks function sent
	select {
	case err := <-errCh:
		assert.Equal(
			t,
			&errTaskExecutorStopped{
				workerID: e.workerID,
				err:      errSome,
			},
			err,
		)
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

func TestRunBlocksUntilReceiveDeferredTasksSendsError(t *testing.T) {
	e := getTestEngine()

	// Override the engine's clean function so it just communicates when the
	// context it was passed has been canceled
	contextCanceledCh := make(chan struct{})
	e.clean = func(
		ctx context.Context,
		_ string,
		_ string,
		_ string,
		_ time.Duration,
	) error {
		<-ctx.Done()
		close(contextCanceledCh)
		return ctx.Err()
	}

	// Override the engine's receiveDeferredTasks function so it just sends an
	// error
	e.receiveDeferredTasks = func(
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
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Call Run in a goroutine. If it never unblocks, as we hope it does, we don't
	// want the test to stall.
	errCh := make(chan error)
	go func() {
		errCh <- e.Run(ctx)
	}()

	// Assert that the error returned from the Run function wraps the error that
	// the overridden receiveDeferredTasks function sent
	select {
	case err := <-errCh:
		assert.Equal(
			t,
			&errReceiverStopped{
				workerID:  e.workerID,
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

func TestRunBlocksUntilWatchDeferredTasksSendsError(t *testing.T) {
	e := getTestEngine()

	// Override the engine's clean function so it just communicates when the
	// context it was passed has been canceled
	contextCanceledCh := make(chan struct{})
	e.clean = func(
		ctx context.Context,
		_ string,
		_ string,
		_ string,
		_ time.Duration,
	) error {
		<-ctx.Done()
		close(contextCanceledCh)
		return ctx.Err()
	}

	// Override the engine's watchDeferredTasks function so it just sends an error
	e.watchDeferredTasks = func(
		ctx context.Context,
		_ chan []byte,
		_ string,
		errCh chan error,
	) {
		select {
		case errCh <- errSome:
		case <-ctx.Done():
		}
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Call Run in a goroutine. If it never unblocks, as we hope it does, we don't
	// want the test to stall.
	errCh := make(chan error)
	go func() {
		errCh <- e.Run(ctx)
	}()

	// Assert that the error returned from the Run function wraps the error that
	// the overridden watchDeferredTask function sent
	select {
	case err := <-errCh:
		assert.Equal(
			t,
			&errDeferredTaskWatcherStopped{workerID: e.workerID, err: errSome},
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

func TestRunRespondsToCanceledContext(t *testing.T) {
	e := getTestEngine()

	// Override the engine's clean function so it just communicates when the
	// context it was passed has been canceled
	contextCanceledCh := make(chan struct{})
	e.clean = func(
		ctx context.Context,
		_ string,
		_ string,
		_ string,
		_ time.Duration,
	) error {
		<-ctx.Done()
		close(contextCanceledCh)
		return ctx.Err()
	}

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

// getTestEngine returns a pointer to an engine that has all its long-running
// concurrent functions pre-overridden to simply block until the context they
// are passed is canceled. Individual test cases can selectively revert or
// amend these overrides to test specific scenarios.
func getTestEngine() *engine {
	config := NewConfigWithDefaults()
	config.RedisHost = "redis"
	config.RedisDB = 1
	config.PendingTaskWorkerCount = 1
	config.DeferedTaskWatcherCount = 1
	e := NewEngine(config).(*engine)
	// Cleaner loop
	e.clean = func(
		ctx context.Context,
		_ string,
		_ string,
		_ string,
		_ time.Duration,
	) error {
		<-ctx.Done()
		return ctx.Err()
	}
	// Heartbeat loop
	e.runHeart = func(ctx context.Context, _ time.Duration) error {
		<-ctx.Done()
		return ctx.Err()
	}
	// Pending tasks receiver
	e.receivePendingTasks = func(
		ctx context.Context,
		_ string,
		_ string,
		_ chan []byte,
		_ chan error,
	) {
		<-ctx.Done()
	}
	// Deferred tasks receiver
	e.receiveDeferredTasks = func(
		ctx context.Context,
		_ string,
		_ string,
		_ chan []byte,
		_ chan error,
	) {
		<-ctx.Done()
	}
	// Tasks executor
	e.executeTasks = func(
		ctx context.Context,
		_ chan []byte,
		_ string,
		_ string,
		_ chan error,
	) {
		<-ctx.Done()
	}
	// Deferred task watcher
	e.watchDeferredTasks = func(
		ctx context.Context,
		_ chan []byte,
		_ string,
		errCh chan error,
	) {
		<-ctx.Done()
	}
	return e
}

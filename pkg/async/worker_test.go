package async

import (
	"context"
	"errors"
	"testing"
	"time"

	fakeAsync "github.com/Azure/open-service-broker-azure/pkg/async/fake"
	"github.com/Azure/open-service-broker-azure/pkg/async/model"
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
	h := fakeAsync.NewHeart()

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
	h := fakeAsync.NewHeart()

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
	h := fakeAsync.NewHeart()

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
	h := fakeAsync.NewHeart()

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
	h := fakeAsync.NewHeart()

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
	h := fakeAsync.NewHeart()

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

// TestDefaultReceiveTasks tests the happy path for the defaultReceiveTasks
// function that transplants tasks from a source queue into a destination queue.
func TestDefaultReceiveTasks(t *testing.T) {
	w := newWorker(redisClient).(*worker)

	sourceQueueName := getDisposableQueueName()
	destinationQueueName := getDisposableQueueName()

	// Put some tasks on the source task queue
	const taskCount int64 = 5
	for range [taskCount]struct{}{} {
		// Dummy tasks are fine. This test won't ever parse them.
		err := redisClient.LPush(sourceQueueName, "foo").Err()
		assert.Nil(t, err)
	}

	// Assert that the source queue has precisely taskCount tasks
	sourceQueueDepth, err := redisClient.LLen(sourceQueueName).Result()
	assert.Nil(t, err)
	assert.Equal(t, taskCount, sourceQueueDepth)

	// Assert that the destination queue is empty
	destinationQueueDepth, err := redisClient.LLen(destinationQueueName).Result()
	assert.Nil(t, err)
	assert.Empty(t, destinationQueueDepth)

	// Under nominal conditions, defaultReceiveTasks blocks until the context it
	// is passed is canceled. Use a context that will cancel itself after 1 second
	// to make defaultReceiveTasks STOP working so we can then examine what it
	// accomplished.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	retCh := make(chan []byte)
	errCh := make(chan error)
	go w.defaultReceiveTasks(
		ctx,
		sourceQueueName,
		destinationQueueName,
		retCh,
		errCh,
	)

	// Start another goroutine to receive and count results
	var resCount int64
	go func() {
		for {
			select {
			case <-retCh:
				resCount++
			case <-ctx.Done():
			}
		}
	}()

	select {
	case <-errCh:
		assert.Fail(t, "should not have received any error, but did")
	case <-ctx.Done():
	}

	// Assert that precisely taskCount tasks were placed onto the return channel
	assert.Equal(t, taskCount, resCount)

	// Assert that the source task queue has been drained
	sourceQueueDepth, err = redisClient.LLen(sourceQueueName).Result()
	assert.Nil(t, err)
	assert.Empty(t, sourceQueueDepth)

	// Assert that the destination queue now has precisely taskCount tasks
	destinationQueueDepth, err = redisClient.LLen(destinationQueueName).Result()
	assert.Nil(t, err)
	assert.Equal(t, taskCount, destinationQueueDepth)
}

func TestDefaultExecuteTasks(t *testing.T) {
	w := newWorker(redisClient).(*worker)

	pendingTaskQueueName := getDisposableQueueName()
	activeTaskQueueName := getActiveTaskQueueName(w.id)

	// Register some jobs with the worker
	var badJobCallCount int
	err := w.RegisterJob(
		"badJob",
		func(_ context.Context, _ map[string]string) error {
			badJobCallCount++
			return errors.New("a deliberate error")
		},
	)
	assert.Nil(t, err)
	var goodJobCallCount int
	err = w.RegisterJob(
		"goodJob",
		func(_ context.Context, _ map[string]string) error {
			goodJobCallCount++
			return nil
		},
	)
	assert.Nil(t, err)

	// Define some tasks
	invalidTaskJSON := []byte("bogus")
	unregisteredTask := model.NewTask("nonExistingJob", map[string]string{})
	unregisteredTaskJSON, err := unregisteredTask.ToJSON()
	assert.Nil(t, err)
	badTask := model.NewTask("badJob", map[string]string{})
	badTaskJSON, err := badTask.ToJSON()
	assert.Nil(t, err)
	goodTask := model.NewTask("goodJob", map[string]string{})
	goodTaskJSON, err := goodTask.ToJSON()
	assert.Nil(t, err)
	assert.Nil(t, err)
	tasks := [][]byte{
		invalidTaskJSON,
		unregisteredTaskJSON,
		badTaskJSON,
		goodTaskJSON,
	}

	// Put all the tasks on the worker's active task queue
	for _, task := range tasks {
		err := redisClient.LPush(activeTaskQueueName, task).Err()
		assert.Nil(t, err)
	}

	// Assert that the pending task queue is empty
	pendingTaskQueueDepth, err := redisClient.LLen(pendingTaskQueueName).Result()
	assert.Nil(t, err)
	assert.Empty(t, pendingTaskQueueDepth)

	// Assert that worker's active task queue has precisely len(tasks) tasks
	activeTaskQueueDepth, err := redisClient.LLen(activeTaskQueueName).Result()
	assert.Nil(t, err)
	assert.Equal(t, int64(len(tasks)), activeTaskQueueDepth)

	// Under nominal conditions, defaultExecuteTasks blocks until the context it
	// is passed is canceled. Use a context that will cancel itself after 1 second
	// to make defaultExecuteTasks STOP working so we can then examine what it
	// accomplished.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// Put all the tasks on the defaultExecuteTasks function's input channel
	inputCh := make(chan []byte)
	go func() {
		for _, task := range tasks {
			select {
			case inputCh <- task:
			case <-ctx.Done():
			}
		}
	}()

	errCh := make(chan error)
	go w.defaultExecuteTasks(ctx, inputCh, pendingTaskQueueName, errCh)

	select {
	case <-errCh:
		assert.Fail(t, "should not have received any error, but did")
	case <-ctx.Done():
	}

	// Assert that the pending task queue has precisely one task-- the
	// unprocessable task, which should have been returned to it
	pendingTaskQueueDepth, err = redisClient.LLen(pendingTaskQueueName).Result()
	assert.Nil(t, err)
	assert.Equal(t, int64(1), pendingTaskQueueDepth)

	// Assert that the worker's active task queue is empty-- in all cases, the
	// tasks should have been removed from this queue
	activeTaskQueueDepth, err = redisClient.LLen(activeTaskQueueName).Result()
	assert.Nil(t, err)
	assert.Empty(t, activeTaskQueueDepth)

	// Assert that the indicated jobs were invoked the appropriate number of times
	assert.Equal(t, 1, badJobCallCount)
	assert.Equal(t, 1, goodJobCallCount)
}

// TestDefaultWatchDeferredTaskWithInvalidTask tests a specific failure
// condition for the defaultWatchDeferredTask function. The expected behavior is
// that the task is discarded and the defaultWatchDeferredTask function adds no
// error to its error channel.
func TestDefaultWatchDeferredTaskWithInvalidTask(t *testing.T) {
	w := newWorker(redisClient).(*worker)

	pendingTaskQueueName := getDisposableQueueName()
	watchedTaskQueueName := getWatchedTaskQueueName(w.id)

	// Assert that the pending task queue is empty
	pendingTaskQueueDepth, err := redisClient.LLen(pendingTaskQueueName).Result()
	assert.Nil(t, err)
	assert.Empty(t, pendingTaskQueueDepth)

	// Put an invalid task (it isn't even JSON) on the worker's watched task queue
	invalidTaskJSON := []byte("bogus")
	err = redisClient.LPush(watchedTaskQueueName, invalidTaskJSON).Err()
	assert.Nil(t, err)

	// Assert that queue now has precisely 1 task
	watchedTaskQueueDepth, err := redisClient.LLen(watchedTaskQueueName).Result()
	assert.Nil(t, err)
	assert.Equal(t, int64(1), watchedTaskQueueDepth)

	// Under nominal conditions, defaultWatchDeferredTask could block for a very
	// long time, unless the context it is passed is canceled. Use a context that
	// will cancel itself after 1 second to make defaultWatchDeferredTask STOP
	// working so we can then examine what it accomplished.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// Call defaultWatchDeferredTask in a goroutine. Under nominal conditions,
	// this function has the potential to run for a long time. In case the
	// function doesn't handle the failure case we're testing properly and return
	// quickly, we do not want to stall this test.
	errCh := make(chan error)
	go w.defaultWatchDeferredTask(
		ctx,
		invalidTaskJSON,
		pendingTaskQueueName,
		errCh,
	)
	select {
	case <-errCh:
		assert.Fail(t, "should not have received any error, but did")
	case <-ctx.Done():
	}

	// Assert that the pending task queue is STILL empty
	pendingTaskQueueDepth, err = redisClient.LLen(pendingTaskQueueName).Result()
	assert.Nil(t, err)
	assert.Empty(t, pendingTaskQueueDepth)

	// And so is the worker's watched task queue
	watchedTaskQueueDepth, err = redisClient.LLen(watchedTaskQueueName).Result()
	assert.Nil(t, err)
	assert.Empty(t, watchedTaskQueueDepth)
}

// TestDefaultWatchDeferredTaskWithTaskWithoutExecuteTime tests a specific
// failure condition for the defaultWatchDeferredTask function. The expected
// behavior is that the task is discarded and the defaultWatchDeferredTask
// function adds no error to its error channel.
func TestDefaultWatchDeferredTaskWithTaskWithoutExecuteTime(t *testing.T) {
	w := newWorker(redisClient).(*worker)

	pendingTaskQueueName := getDisposableQueueName()
	watchedTaskQueueName := getWatchedTaskQueueName(w.id)

	// Assert that the pending task queue is empty
	pendingTaskQueueDepth, err := redisClient.LLen(pendingTaskQueueName).Result()
	assert.Nil(t, err)
	assert.Empty(t, pendingTaskQueueDepth)

	// Put a task with no execute time on the worker's watched task queue
	task := model.NewTask("foo", nil)
	taskJSON, err := task.ToJSON()
	assert.Nil(t, err)
	err = redisClient.LPush(watchedTaskQueueName, taskJSON).Err()
	assert.Nil(t, err)

	// Assert that queue now has precisely 1 task
	watchedTaskQueueDepth, err := redisClient.LLen(watchedTaskQueueName).Result()
	assert.Nil(t, err)
	assert.Equal(t, int64(1), watchedTaskQueueDepth)

	// Under nominal conditions, defaultWatchDeferredTask could block for a very
	// long time, unless the context it is passed is canceled. Use a context that
	// will cancel itself after 1 second to make defaultWatchDeferredTask STOP
	// working so we can then examine what it accomplished.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// Call defaultWatchDeferredTask in a goroutine. Under nominal conditions,
	// this function has the potential to run for a long time. In case the
	// function doesn't handle the failure case we're testing properly and return
	// quickly, we do not want to stall this test.
	errCh := make(chan error)
	go w.defaultWatchDeferredTask(
		ctx,
		taskJSON,
		pendingTaskQueueName,
		errCh,
	)
	select {
	case <-errCh:
		assert.Fail(t, "should not have received any error, but did")
	case <-ctx.Done():
	}

	// Assert that the pending task queue is STILL empty
	pendingTaskQueueDepth, err = redisClient.LLen(pendingTaskQueueName).Result()
	assert.Nil(t, err)
	assert.Empty(t, pendingTaskQueueDepth)

	// And so is the worker's watched task queue
	watchedTaskQueueDepth, err = redisClient.LLen(watchedTaskQueueName).Result()
	assert.Nil(t, err)
	assert.Empty(t, watchedTaskQueueDepth)
}

// TestDefaultWatchDeferredTaskWithLapsedTask tests that when
// defaultWatchDeferredTask is invoked for a task whose execute time has already
// lapsed, that task is moved IMMEDIATELY to the pending task queue.
func TestDefaultWatchDeferredTaskWithLapsedTask(t *testing.T) {
	w := newWorker(redisClient).(*worker)

	pendingTaskQueueName := getDisposableQueueName()
	watchedTaskQueueName := getWatchedTaskQueueName(w.id)

	// Assert that the pending task queue is empty
	pendingTaskQueueDepth, err := redisClient.LLen(pendingTaskQueueName).Result()
	assert.Nil(t, err)
	assert.Empty(t, pendingTaskQueueDepth)

	// Put a lapsed task on the worker's watched task queue
	task := model.NewDelayedTask("foo", nil, time.Second*-1)
	taskJSON, err := task.ToJSON()
	assert.Nil(t, err)
	err = redisClient.LPush(watchedTaskQueueName, taskJSON).Err()
	assert.Nil(t, err)

	// Assert that queue now has precisely 1 task
	watchedTaskQueueDepth, err := redisClient.LLen(watchedTaskQueueName).Result()
	assert.Nil(t, err)
	assert.Equal(t, int64(1), watchedTaskQueueDepth)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Call defaultWatchDeferredTask in a goroutine. Under nominal conditions,
	// this function has the potential to run for a long time, although in this
	// edge case, it should not. In case the function doesn't handle the edge case
	// we're testing properly and return quickly, we do not want to stall this
	// test.
	errCh := make(chan error)
	go w.defaultWatchDeferredTask(
		ctx,
		taskJSON,
		pendingTaskQueueName,
		errCh,
	)
	select {
	case <-errCh:
		assert.Fail(t, "should not have received any error, but did")
	case <-time.After(time.Second):
	}

	// Assert that the pending task queue now has precisely 1 task
	pendingTaskQueueDepth, err = redisClient.LLen(pendingTaskQueueName).Result()
	assert.Nil(t, err)
	assert.Equal(t, int64(1), pendingTaskQueueDepth)

	// And the worker's watched task queue is now empty
	watchedTaskQueueDepth, err = redisClient.LLen(watchedTaskQueueName).Result()
	assert.Nil(t, err)
	assert.Empty(t, watchedTaskQueueDepth)
}

// TestDefaultWatchDeferredTaskRespondsToCanceledContext tests that when
// defaultWatchDeferredTask is waiting for a tasks execute time to lapse, it
// will abort if context is canceled.
func TestDefaultWatchDeferredTaskRespondsToCanceledContext(t *testing.T) {
	w := newWorker(redisClient).(*worker)

	pendingTaskQueueName := getDisposableQueueName()
	watchedTaskQueueName := getWatchedTaskQueueName(w.id)

	// Assert that the pending task queue is empty
	pendingTaskQueueDepth, err := redisClient.LLen(pendingTaskQueueName).Result()
	assert.Nil(t, err)
	assert.Empty(t, pendingTaskQueueDepth)

	// Put task with a future execute time on the worker's watched tasks queue
	task := model.NewDelayedTask("foo", nil, time.Second*5)
	taskJSON, err := task.ToJSON()
	assert.Nil(t, err)
	err = redisClient.LPush(watchedTaskQueueName, taskJSON).Err()
	assert.Nil(t, err)

	// Assert that queue now has precisely 1 task
	watchedTaskQueueDepth, err := redisClient.LLen(watchedTaskQueueName).Result()
	assert.Nil(t, err)
	assert.Equal(t, int64(1), watchedTaskQueueDepth)

	// Use a context that will cancel itself in 1 second to put a time limit
	// on the test. This context will be canceled BEFORE the execute time lapses.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// Call defaultWatchDeferredTask in a goroutine. Under nominal conditions,
	// this function has the potential to run for a long time. In case the
	// function doesn't handle context cancelation properly and return quickly, we
	// do not want to stall this test.
	errCh := make(chan error)
	go w.defaultWatchDeferredTask(
		ctx,
		taskJSON,
		pendingTaskQueueName,
		errCh,
	)
	select {
	case <-errCh:
		assert.Fail(t, "should not have received any error, but did")
	case <-ctx.Done():
	}

	// Because the context was canceled BEFORE the execute time lapsed, nothing
	// should have changed in the queues....

	// Assert that the pending task queue is STILL empty
	pendingTaskQueueDepth, err = redisClient.LLen(pendingTaskQueueName).Result()
	assert.Nil(t, err)
	assert.Empty(t, pendingTaskQueueDepth)

	// And the worker's watched task STILL has precisely 1 task
	watchedTaskQueueDepth, err = redisClient.LLen(watchedTaskQueueName).Result()
	assert.Nil(t, err)
	assert.Equal(t, int64(1), watchedTaskQueueDepth)
}

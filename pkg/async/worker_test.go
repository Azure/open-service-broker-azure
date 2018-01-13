package async

import (
	"context"
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

// TestWorkBlocksUntilHeartReturns tests what happens when a worker's heart
// stops beating. The expected behavior is that everything else the worker is
// doing concurrently has its context canceled and the Work function returns an
// appropriate errHeartStopped error.
func TestWorkBlocksUntilHeartReturns(t *testing.T) {
	// Use a fake heart
	h := fakeAsync.NewHeart()

	// Specify the fake heart's runtime behavior should just return an error
	h.RunBehavior = func(context.Context) error {
		return errSome
	}

	w := newWorker(redisClient).(*worker)

	// Make the worker use the fake heart
	w.heart = h

	// Override the worker's default receiveAndWork function so it just
	// communicates when the context it was passed has been canceled
	contextCanceledCh := make(chan struct{})
	w.receiveAndWork = func(ctx context.Context, _ string) error {
		<-ctx.Done()
		select {
		case contextCanceledCh <- struct{}{}:
		default:
		}
		return ctx.Err()
	}

	// Override the worker's default handleDelayedTasks function so it just blocks
	// until the context it was passed is canceled
	w.handleDelayedTasks = func(ctx context.Context, _ string, _ string) error {
		<-ctx.Done()
		return ctx.Err()
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	errCh := make(chan error)
	go func() {
		errCh <- w.Work(ctx)
	}()

	// Assert that the error returned from the Work function wraps the error that
	// the fake heart generated
	select {
	case err := <-errCh:
		assert.Equal(t, &errHeartStopped{workerID: w.id, err: errSome}, err)
	case <-time.After(time.Second):
		assert.Fail(
			t,
			"an errHeartStopped error should have been returned, but wasn't",
		)
	}

	// Assert that the context got canceled
	select {
	case <-contextCanceledCh:
	case <-time.After(time.Second):
		assert.Fail(t, "context should have been canceled, but it was not")
	}
}

// TestWorkBlocksUntilReceiveAndWorkReturns tests what happens when a worker's
// main loop returns an error. The expected behavior is that everything else the
// worker is doing concurrently has its context canceled and the Work function
// returns an appropriate errReceiveAndWorkStopped error.
func TestWorkBlocksUntilReceiveAndWorkReturns(t *testing.T) {
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

	// Override the worker's default receiveAndWork function so it just returns an
	// error
	w.receiveAndWork = func(context.Context, string) error {
		return errSome
	}

	// Override the worker's default handleDelayedTasks function so it just blocks
	// until its context is canceled
	w.handleDelayedTasks = func(ctx context.Context, _ string, _ string) error {
		<-ctx.Done()
		return ctx.Err()
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	errCh := make(chan error)
	go func() {
		errCh <- w.Work(ctx)
	}()

	// Assert that the error received from the Work function wraps the error that
	// the overridden receiveAndWork function generated
	select {
	case err := <-errCh:
		assert.Equal(
			t,
			&errReceiveAndWorkStopped{workerID: w.id, err: errSome},
			err,
		)
	case <-time.After(time.Second):
		assert.Fail(
			t,
			"an errReceiveAndWorkStopped error should have been returned, but wasn't",
		)
	}

	// Assert that the context got canceled
	select {
	case <-contextCanceledCh:
	case <-time.After(time.Second):
		assert.Fail(t, "context should have been canceled, but it was not")
	}
}

// TestWorkBlocksUntilHandleDelayedTasksReturns tests what happens when a
// worker's delayed task handler returns an error. The expected behavior is that
// everything else the worker is doing concurrently has its context canceled and
// the Work function returns an appropriate errWatchDelayedTasksStopped error.
func TestWorkBlocksUntilHandleDelayedTasksReturns(t *testing.T) {
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

	// Override the worker's default receiveAndWork function so it just blocks
	// until its context is canceled
	w.receiveAndWork = func(ctx context.Context, _ string) error {
		<-ctx.Done()
		return ctx.Err()
	}

	// Override the worker's default handleDelayedTasks function so it just
	// returns an error
	w.handleDelayedTasks = func(context.Context, string, string) error {
		return errSome
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	errCh := make(chan error)
	go func() {
		errCh <- w.Work(ctx)
	}()

	// Assert that the error received from the Work function wraps the error that
	// the overridden handleDelayedTasks function generated
	select {
	case err := <-errCh:
		assert.Equal(
			t,
			&errWatchDelayedTasksStopped{workerID: w.id, err: errSome},
			err,
		)
	case <-time.After(time.Second):
		assert.Fail(
			t,
			"an errWatchDelayedTasksStopped error should have been returned, but wasn't",
		)
	}

	// Assert that the context got canceled
	select {
	case <-contextCanceledCh:
	case <-time.After(time.Second):
		assert.Fail(t, "context should have been canceled, but it was not")
	}
}

// TestWorkBlocksUntilContextCanceled tests that canceling the context passed
// to the Work function causes the Work function to return.
func TestWorkBlocksUntilContextCanceled(t *testing.T) {
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

	// Override the worker's default receiveAndWork function so it just blocks
	// until its context is canceled
	w.receiveAndWork = func(ctx context.Context, _ string) error {
		<-ctx.Done()
		return ctx.Err()
	}

	// Override the worker's default handleDelayedTasks function so it just blocks
	// until its context is canceled
	w.handleDelayedTasks = func(ctx context.Context, _ string, _ string) error {
		<-ctx.Done()
		return ctx.Err()
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	errCh := make(chan error)
	go func() {
		errCh <- w.Work(ctx)
	}()

	cancel()

	// Assert that the error returned from Work indicates that the context was
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

// TestReceiveAndWork tests the happy path for the receiveAndWork function that
// handles tasks from the main active work queue. The expected behavior is that
// another internal work function is invoked once per task and that the main
// active work queue is drained. The worker's own active work queue also should
// be empty when this test completes.
func TestReceiveAndWork(t *testing.T) {
	w := newWorker(redisClient).(*worker)

	// Override the worker's default work function so it just counts how many
	// times it has been invoked
	var workCallCount int64
	w.work = func(context.Context, model.Task) error {
		workCallCount++
		return nil
	}

	mainActiveWorkQueueName := getDisposableQueueName()
	workerActiveWorkQueueName := getWorkerActiveQueueName(w.id)

	// Put some tasks on the main active work queue.
	const taskCount int64 = 5
	for range [taskCount]struct{}{} {
		taskJSON, err := model.NewTask("foo", nil).ToJSON()
		assert.Nil(t, err)
		intCmd := redisClient.LPush(mainActiveWorkQueueName, taskJSON)
		assert.Nil(t, intCmd.Err())
	}

	// Assert that the main active work queue has precisely taskCount tasks
	intCmd := redisClient.LLen(mainActiveWorkQueueName)
	assert.Nil(t, intCmd.Err())
	mainActiveWorkQueueDepth, err := intCmd.Result()
	assert.Nil(t, err)
	assert.Equal(t, taskCount, mainActiveWorkQueueDepth)

	// Assert that the worker's active work queue is empty
	intCmd = redisClient.LLen(workerActiveWorkQueueName)
	assert.Nil(t, intCmd.Err())
	workerActiveWorkQueueDepth, err := intCmd.Result()
	assert.Nil(t, err)
	assert.Empty(t, workerActiveWorkQueueDepth)

	// Use a context that will cancel itself in 1 second because we want the
	// receiveAndWork function to return soon instead of blocking indefinitely
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	errCh := make(chan error)
	go func() {
		errCh <- w.receiveAndWork(ctx, mainActiveWorkQueueName)
	}()

	// Assert that the error we received is only an indicator of context
	// cancelation
	select {
	case err := <-errCh:
		assert.Equal(t, ctx.Err(), err)
	// Note that we do blocking reads from the main active work queue that take
	// 5 seconds to time out. So we only have the opporunity to deal with a
	// canceled context every 5 seconds. Setting the timeout for this select to
	// 6 seconds should give us enough time to process the 5 tasks on the main
	// active work queue and time out on the 6th read, at which point the canceled
	// context should be dealt with.
	case <-time.After(time.Second * 6):
		assert.Fail(
			t,
			"a context canceled error should have been returned, but wasn't",
		)
	}

	// Assert that the internal work function was called exactly once for each
	// task we placed on the main active work queue
	assert.Equal(t, taskCount, workCallCount)

	// Assert that the main active work queue has been drained
	intCmd = redisClient.LLen(mainActiveWorkQueueName)
	assert.Nil(t, intCmd.Err())
	mainActiveWorkQueueDepth, err = intCmd.Result()
	assert.Nil(t, err)
	assert.Empty(t, mainActiveWorkQueueDepth)

	// Assert that the worker's active work queue is also empty
	intCmd = redisClient.LLen(workerActiveWorkQueueName)
	assert.Nil(t, intCmd.Err())
	workerActiveWorkQueueDepth, err = intCmd.Result()
	assert.Nil(t, err)
	assert.Empty(t, workerActiveWorkQueueDepth)
}

// TestReceiveAndWorkKeepsRunningAfterInvalidTask tests a specific failure
// condition for the receiveAndWork function that handles tasks from the main
// active work queue. The expected behavior is that an invalid task on the
// main active work queue is simply discarded and does not cause the
// receiveAndWork function to stop handling tasks.
func TestReceiveAndWorkKeepsRunningAfterInvalidTask(t *testing.T) {
	w := newWorker(redisClient).(*worker)

	// Override the worker's default work function so it just records that it's
	// been invoked
	workCalled := false
	w.work = func(context.Context, model.Task) error {
		workCalled = true
		return nil
	}

	mainActiveWorkQueueName := getDisposableQueueName()
	workerActiveWorkQueueName := getWorkerActiveQueueName(w.id)

	// Put an invalid task (it isn't even JSON) on the main active work queue
	intCmd := redisClient.LPush(mainActiveWorkQueueName, "bogus")
	assert.Nil(t, intCmd.Err())

	// Assert that the main active work queue has precisely 1 task
	intCmd = redisClient.LLen(mainActiveWorkQueueName)
	assert.Nil(t, intCmd.Err())
	mainActiveWorkQueueDepth, err := intCmd.Result()
	assert.Nil(t, err)
	assert.Equal(t, int64(1), mainActiveWorkQueueDepth)

	// Assert that the worker's active work queue is empty
	intCmd = redisClient.LLen(workerActiveWorkQueueName)
	assert.Nil(t, intCmd.Err())
	workerActiveWorkQueueDepth, err := intCmd.Result()
	assert.Nil(t, err)
	assert.Empty(t, workerActiveWorkQueueDepth)

	// Use a context that will cancel itself in 1 second because we want the
	// receiveAndWork function to return soon instead of blocking indefinitely
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	errCh := make(chan error)
	go func() {
		errCh <- w.receiveAndWork(ctx, mainActiveWorkQueueName)
	}()

	// Assert that the error we received is ONLY an indicator of context
	// cancelation. i.e. We didn't fail because of this.
	select {
	case err := <-errCh:
		assert.Equal(t, ctx.Err(), err)
	// Note that we do blocking reads from the main active work queue that take
	// 5 seconds to time out. So we only have the opporunity to deal with a
	// canceled context every 5 seconds. Setting the timeout for this select to
	// 6 seconds should give us enough time to process the 1 invalid task on the
	// main active work queue and time out on the 2nd read, at which point the
	// canceled context should be dealt with.
	case <-time.After(time.Second * 6):
		assert.Fail(
			t,
			"a context canceled error should have been returned, but wasn't",
		)
	}

	// Assert that the work function was never invoked
	assert.False(t, workCalled)

	// Assert that the main active work queue has been drained
	intCmd = redisClient.LLen(mainActiveWorkQueueName)
	assert.Nil(t, intCmd.Err())
	mainActiveWorkQueueDepth, err = intCmd.Result()
	assert.Nil(t, err)
	assert.Empty(t, mainActiveWorkQueueDepth)

	// Assert that the worker's active work queue is also empty
	intCmd = redisClient.LLen(workerActiveWorkQueueName)
	assert.Nil(t, intCmd.Err())
	workerActiveWorkQueueDepth, err = intCmd.Result()
	assert.Nil(t, err)
	assert.Empty(t, workerActiveWorkQueueDepth)
}

// TestReceiveAndWorkKeepsRunningAfterWorkError tests a specific failure
// condition for the receiveAndWork function that handles tasks from the main
// active work queue. The expected behavior is that an error executing a task
// does not cause the receiveAndWork function to stop handling tasks.
func TestReceiveAndWorkKeepsRunningAfterWorkError(t *testing.T) {
	w := newWorker(redisClient).(*worker)

	// Override the worker's work function so it records that it's been invoked
	// and returns an error
	workCalled := false
	w.work = func(context.Context, model.Task) error {
		workCalled = true
		return errSome
	}

	mainActiveWorkQueueName := getDisposableQueueName()
	workerActiveWorkQueueName := getWorkerActiveQueueName(w.id)

	// Use a valid task because it will be parsed in the course of this test
	taskJSON, err := model.NewTask("foo", nil).ToJSON()
	assert.Nil(t, err)

	// Put the task on the main active work queue
	intCmd := redisClient.LPush(mainActiveWorkQueueName, taskJSON)
	assert.Nil(t, intCmd.Err())

	// Assert the main active work queue has precisely 1 task
	intCmd = redisClient.LLen(mainActiveWorkQueueName)
	assert.Nil(t, intCmd.Err())
	mainActiveWorkQueueDepth, err := intCmd.Result()
	assert.Nil(t, err)
	assert.Equal(t, int64(1), mainActiveWorkQueueDepth)

	// Assert that the worker's active work queue is empty
	intCmd = redisClient.LLen(workerActiveWorkQueueName)
	assert.Nil(t, intCmd.Err())
	workerActiveWorkQueueDepth, err := intCmd.Result()
	assert.Nil(t, err)
	assert.Empty(t, workerActiveWorkQueueDepth)

	// Use a context that will cancel itself in 1 second because we want the
	// receiveAndWork function to return soon instead of blocking indefinitely
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
	defer cancel()

	errCh := make(chan error)
	go func() {
		errCh <- w.receiveAndWork(ctx, mainActiveWorkQueueName)
	}()

	// Assert that the error we received is ONLY an indicator of context
	// cancelation. i.e. We didn't fail because of this.
	select {
	case err := <-errCh:
		assert.Equal(t, ctx.Err(), err)
	// Note that we do blocking reads from the main active work queue that take
	// 5 seconds to time out. So we only have the opporunity to deal with a
	// canceled context every 5 seconds. Setting the timeout for this select to
	// 6 seconds should give us enough time to process the 1 task on the main
	// active work queue (which will faile) and then time out on the 2nd read, at
	// which point the canceled context should be dealt with.
	case <-time.After(time.Second * 6):
		assert.Fail(
			t,
			"a context canceled error should have been returned, but wasn't",
		)
	}

	// Assert that the work function was invoked
	assert.True(t, workCalled)

	// Assert that the main active work queue has been drained
	intCmd = redisClient.LLen(mainActiveWorkQueueName)
	assert.Nil(t, intCmd.Err())
	mainActiveWorkQueueDepth, err = intCmd.Result()
	assert.Nil(t, err)
	assert.Empty(t, mainActiveWorkQueueDepth)

	// Assert that the worker's active work queue is also empty
	intCmd = redisClient.LLen(workerActiveWorkQueueName)
	assert.Nil(t, intCmd.Err())
	workerActiveWorkQueueDepth, err = intCmd.Result()
	assert.Nil(t, err)
	assert.Empty(t, workerActiveWorkQueueDepth)
}

// TestReceiveAndWorkBlocksUntilContextCanceled tests that canceling the context
// passed to the receiveAndWork function causes the receiveAndWork function to
// return.
func TestReceiveAndWorkBlocksUntilContextCanceled(t *testing.T) {
	w := newWorker(redisClient).(*worker)

	// Use a context that will cancel itself in 1 second because we want the
	// receiveAndWork function to return soon instead of blocking indefinitely
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	errCh := make(chan error)
	go func() {
		errCh <- w.receiveAndWork(ctx, getDisposableQueueName())
	}()

	// Assert that the error returned from receiveAndWork indicates that the
	// context was canceled
	select {
	case err := <-errCh:
		assert.Equal(t, ctx.Err(), err)
	// Note that we do blocking reads from the main active work queue that take
	// 5 seconds to time out. So we only have the opporunity to deal with a
	// canceled context every 5 seconds. Setting the timeout for this select to
	// 6 seconds should give us enough time to time out on the 1st read, at
	// which point the canceled context should be dealt with.
	case <-time.After(time.Second * 6):
		assert.Fail(
			t,
			"a context canceled error should have been returned, but wasn't",
		)
	}
}

// TestHandleDelayedTasks tests the happy path for the handleDelayedTasks
// function that handles tasks from the main delayed work queue. The expected
// behavior is that handleDelayedTask is launched in a new goroutine once per
// task and that the main delayed work queue is drained.
func TestHandleDelayedTasks(t *testing.T) {
	w := newWorker(redisClient).(*worker)

	// Override the worker's handleDelayedTask function so it just counts how many
	// times it has been invoked
	var handleDelayedTaskCallCount int64
	w.handleDelayedTask = func(
		context.Context,
		[]byte,
		string,
		chan error,
	) {
		handleDelayedTaskCallCount++
	}

	mainActiveWorkQueueName := getDisposableQueueName()
	mainDelayedWorkQueueName := getDisposableQueueName()
	workerDelayedWorkQueueName := getWorkerDelayedQueueName(w.id)

	// Put some tasks on the main delayed work queue.
	const taskCount int64 = 5
	for range [taskCount]struct{}{} {
		// Dummy tasks are fine. This test won't ever parse them.
		intCmd := redisClient.LPush(mainDelayedWorkQueueName, "foo")
		assert.Nil(t, intCmd.Err())
	}

	// Assert that the main delayed work queue has precisely taskCount tasks
	intCmd := redisClient.LLen(mainDelayedWorkQueueName)
	assert.Nil(t, intCmd.Err())
	mainDelayedWorkQueueDepth, err := intCmd.Result()
	assert.Nil(t, err)
	assert.Equal(t, taskCount, mainDelayedWorkQueueDepth)

	// Assert that the worker's delayed work queue is empty
	intCmd = redisClient.LLen(workerDelayedWorkQueueName)
	assert.Nil(t, intCmd.Err())
	workerDelayedWorkQueueDepth, err := intCmd.Result()
	assert.Nil(t, err)
	assert.Empty(t, workerDelayedWorkQueueDepth)

	// Use a context that will cancel itself in 1 second because we want the
	// handleDelayedTasks function to return soon instead of blocking indefinitely
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	errCh := make(chan error)
	go func() {
		errCh <- w.handleDelayedTasks(
			ctx,
			mainDelayedWorkQueueName,
			mainActiveWorkQueueName,
		)
	}()

	// Assert that the error we received is only an indicator of context
	// cancelation
	select {
	case err := <-errCh:
		assert.Equal(t, ctx.Err(), err)
	// Note that we do blocking reads from the main delayed work queue that take
	// 5 seconds to time out. So we only have the opporunity to deal with a
	// canceled context every 5 seconds. Setting the timeout for this select to
	// 6 seconds should give us enough time to process the 5 tasks on the delayed
	// active work queue and time out on the 6th read, at which point the canceled
	// context should be dealt with.
	case <-time.After(time.Second * 6):
		assert.Fail(
			t,
			"a context canceled error should have been returned, but wasn't",
		)
	}

	// Assert that the handleDelayedTask function was called exactly once for each
	// dummy task we placed on the main delayed work queue
	assert.Equal(t, taskCount, handleDelayedTaskCallCount)

	// Assert that the main delayed work queue has been drained
	intCmd = redisClient.LLen(mainDelayedWorkQueueName)
	assert.Nil(t, intCmd.Err())
	mainDelayedWorkQueueDepth, err = intCmd.Result()
	assert.Nil(t, err)
	assert.Empty(t, mainDelayedWorkQueueDepth)

	// Assert that the worker's delayed work queue is taskCount deep
	intCmd = redisClient.LLen(workerDelayedWorkQueueName)
	assert.Nil(t, intCmd.Err())
	workerDelayedWorkQueueDepth, err = intCmd.Result()
	assert.Nil(t, err)
	assert.Equal(t, taskCount, workerDelayedWorkQueueDepth)
}

// TestHandleDelayedTaskWithInvalidTask tests a specific failure condition
// for the handleDelayedTask function that handles delayed tasks from the
// worker's delayed work queue. The expected behavior is that the task is
// discarded and the handleDelayedTask adds no error to its error channel.
func TestHandleDelayedTaskWithInvalidTask(t *testing.T) {
	w := newWorker(redisClient).(*worker)

	mainActiveWorkQueueName := getDisposableQueueName()
	workerDelayedWorkQueueName := getWorkerDelayedQueueName(w.id)

	// Assert that the main active work queue is empty
	intCmd := redisClient.LLen(mainActiveWorkQueueName)
	assert.Nil(t, intCmd.Err())
	mainActiveWorkQueueDepth, err := intCmd.Result()
	assert.Nil(t, err)
	assert.Empty(t, mainActiveWorkQueueDepth)

	// Put an invalid task (it isn't even JSON) on the worker's delayed work queue
	intCmd = redisClient.LPush(workerDelayedWorkQueueName, "bogus")
	assert.Nil(t, intCmd.Err())

	// Assert that queue now has precisely 1 task
	intCmd = redisClient.LLen(workerDelayedWorkQueueName)
	assert.Nil(t, intCmd.Err())
	workerDelayedQueueDepth, err := intCmd.Result()
	assert.Nil(t, err)
	assert.Equal(t, int64(1), workerDelayedQueueDepth)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	errCh := make(chan error)
	go w.handleDelayedTask(ctx, []byte("bogus"), mainActiveWorkQueueName, errCh)
	select {
	case <-errCh:
		assert.Fail(t, "should not have received any error, but did")
	case <-time.After(time.Second):
	}

	// Assert that the main active work queue is STILL empty
	intCmd = redisClient.LLen(mainActiveWorkQueueName)
	assert.Nil(t, intCmd.Err())
	mainActiveWorkQueueDepth, err = intCmd.Result()
	assert.Nil(t, err)
	assert.Empty(t, mainActiveWorkQueueDepth)

	// And so is the worker's delayed work queue
	intCmd = redisClient.LLen(workerDelayedWorkQueueName)
	assert.Nil(t, intCmd.Err())
	workerDelayedWorkQueueDepth, err := intCmd.Result()
	assert.Nil(t, err)
	assert.Empty(t, workerDelayedWorkQueueDepth)
}

// TestHandleDelayedTaskWithoutExecuteTime tests a specific failure condition
// for the handleDelayedTask function that handles delayed tasks from the
// worker's delayed work queue. The expected behavior is that the task is
// discarded and the handleDelayedTask adds no error to its error channel.
func TestHandleDelayedTaskWithoutExecuteTime(t *testing.T) {
	w := newWorker(redisClient).(*worker)

	mainActiveWorkQueueName := getDisposableQueueName()
	workerDelayedWorkQueueName := getWorkerDelayedQueueName(w.id)

	// Assert that the main active work queue is empty
	intCmd := redisClient.LLen(mainActiveWorkQueueName)
	assert.Nil(t, intCmd.Err())
	mainActiveWorkQueueDepth, err := intCmd.Result()
	assert.Nil(t, err)
	assert.Empty(t, mainActiveWorkQueueDepth)

	// Put a task lacking an execute time on the worker's delayed work queue
	taskJSON, err := model.NewTask("foo", nil).ToJSON()
	assert.Nil(t, err)
	intCmd = redisClient.LPush(workerDelayedWorkQueueName, taskJSON)
	assert.Nil(t, intCmd.Err())

	// Assert that queue now has precisely 1 task
	intCmd = redisClient.LLen(workerDelayedWorkQueueName)
	assert.Nil(t, intCmd.Err())
	workerDelayedQueueDepth, err := intCmd.Result()
	assert.Nil(t, err)
	assert.Equal(t, int64(1), workerDelayedQueueDepth)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	errCh := make(chan error)
	go w.handleDelayedTask(ctx, taskJSON, mainActiveWorkQueueName, errCh)
	select {
	case <-errCh:
		assert.Fail(t, "should not have received any error, but did")
	case <-time.After(time.Second):
	}

	// Assert that the main active work queue is STILL empty
	intCmd = redisClient.LLen(mainActiveWorkQueueName)
	assert.Nil(t, intCmd.Err())
	mainActiveWorkQueueDepth, err = intCmd.Result()
	assert.Nil(t, err)
	assert.Empty(t, mainActiveWorkQueueDepth)

	// And so is the worker's delayed work queue
	intCmd = redisClient.LLen(workerDelayedWorkQueueName)
	assert.Nil(t, intCmd.Err())
	workerDelayedWorkQueueDepth, err := intCmd.Result()
	assert.Nil(t, err)
	assert.Empty(t, workerDelayedWorkQueueDepth)
}

// TestHandleDelayedTaskRespondsToContextCanceled tests that canceling the
// context passed to handleDelayedTask causes that function to return. When this
// happens, the depth of both the main active work queue and the worker's
// delayed work queue should remain unaffected.
func TestHandleDelayedTaskRespondsToContextCanceled(t *testing.T) {
	w := newWorker(redisClient).(*worker)

	mainActiveWorkQueueName := getDisposableQueueName()
	workerDelayedWorkQueue := getWorkerDelayedQueueName(w.id)

	// Assert that the main active work queue is empty
	intCmd := redisClient.LLen(mainActiveWorkQueueName)
	assert.Nil(t, intCmd.Err())
	mainActiveWorkQueueDepth, err := intCmd.Result()
	assert.Nil(t, err)
	assert.Empty(t, mainActiveWorkQueueDepth)

	// Need to use a legitimate delayed task for this test
	task := model.NewDelayedTask("foo", nil, time.Minute)
	taskJSON, err := task.ToJSON()
	assert.Nil(t, err)

	// Put the task on the worker's delayed work queue
	intCmd = redisClient.LPush(workerDelayedWorkQueue, taskJSON)
	assert.Nil(t, intCmd.Err())

	// Assert that queue now has precisely 1 task
	intCmd = redisClient.LLen(workerDelayedWorkQueue)
	assert.Nil(t, intCmd.Err())
	workerDelayedQueueDepth, err := intCmd.Result()
	assert.Nil(t, err)
	assert.Equal(t, int64(1), workerDelayedQueueDepth)

	// Use a context that will cancel itself in 1 second
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	errCh := make(chan error)
	go w.handleDelayedTask(ctx, taskJSON, mainActiveWorkQueueName, errCh)
	select {
	case err := <-errCh:
		assert.Equal(t, ctx.Err(), err)
	case <-time.After(time.Second * 2):
		assert.Fail(
			t,
			"a context canceled error should have been returned, but wasn't",
		)
	}

	// Assert that the main active work queue is STILL empty
	intCmd = redisClient.LLen(mainActiveWorkQueueName)
	assert.Nil(t, intCmd.Err())
	mainActiveWorkQueueDepth, err = intCmd.Result()
	assert.Nil(t, err)
	assert.Empty(t, mainActiveWorkQueueDepth)

	// And assert the worker's delayed work queue STILL has the task in it
	intCmd = redisClient.LLen(workerDelayedWorkQueue)
	assert.Nil(t, intCmd.Err())
	workerDelayedQueueDepth, err = intCmd.Result()
	assert.Nil(t, err)
	assert.Equal(t, int64(1), workerDelayedQueueDepth)
}

// TestHandleDelayedTaskWithLapsedTask tests that when handleDelayedTask is
// invoked for a task whose executeTime has already lapsed, that task is moved
// IMMEDIATELY to the main active work queue.
func TestHandleDelayedTaskWithLapsedTask(t *testing.T) {
	w := newWorker(redisClient).(*worker)

	mainActiveWorkQueueName := getDisposableQueueName()
	workerDelayedWorkQueue := getWorkerDelayedQueueName(w.id)

	// Assert that the main active work queue is empty
	intCmd := redisClient.LLen(mainActiveWorkQueueName)
	assert.Nil(t, intCmd.Err())
	mainActiveWorkQueueDepth, err := intCmd.Result()
	assert.Nil(t, err)
	assert.Empty(t, mainActiveWorkQueueDepth)

	// Need to use a lapsed delayed task for this test
	task := model.NewDelayedTask("foo", nil, time.Minute*-1)
	taskJSON, err := task.ToJSON()
	assert.Nil(t, err)

	// Put the lapsed task on the worker's delayed work queue
	intCmd = redisClient.LPush(workerDelayedWorkQueue, taskJSON)
	assert.Nil(t, intCmd.Err())

	// Assert that queue now has precisely 1 task
	intCmd = redisClient.LLen(workerDelayedWorkQueue)
	assert.Nil(t, intCmd.Err())
	workerDelayedQueueDepth, err := intCmd.Result()
	assert.Nil(t, err)
	assert.Equal(t, int64(1), workerDelayedQueueDepth)

	// Use a context that will cancel itself in 1 second
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	errCh := make(chan error)
	go w.handleDelayedTask(ctx, taskJSON, mainActiveWorkQueueName, errCh)
	select {
	case <-errCh:
		assert.Fail(t, "should not have received any error, but did")
	case <-time.After(time.Second * 2):
	}

	// Assert that the main active work queue has precisely 1 task
	intCmd = redisClient.LLen(mainActiveWorkQueueName)
	assert.Nil(t, intCmd.Err())
	mainActiveWorkQueueDepth, err = intCmd.Result()
	assert.Nil(t, err)
	assert.Equal(t, int64(1), mainActiveWorkQueueDepth)

	// And assert the worker's delayed work queue is now empty
	intCmd = redisClient.LLen(workerDelayedWorkQueue)
	assert.Nil(t, intCmd.Err())
	workerDelayedQueueDepth, err = intCmd.Result()
	assert.Nil(t, err)
	assert.Empty(t, workerDelayedQueueDepth)
}

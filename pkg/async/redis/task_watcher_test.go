package redis

import (
	"context"
	"testing"
	"time"

	"github.com/Azure/open-service-broker-azure/pkg/async"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

// TestDefaultWatchDeferredTaskWithInvalidTask tests a specific failure
// condition for the defaultWatchDeferredTask function. The expected behavior is
// that the task is discarded and the defaultWatchDeferredTask function adds no
// error to its error channel.
func TestDefaultWatchDeferredTaskWithInvalidTask(t *testing.T) {
	w := newWorker(redisClient, uuid.NewV4().String()).(*worker)

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
	w := newWorker(redisClient, uuid.NewV4().String()).(*worker)

	pendingTaskQueueName := getDisposableQueueName()
	watchedTaskQueueName := getWatchedTaskQueueName(w.id)

	// Assert that the pending task queue is empty
	pendingTaskQueueDepth, err := redisClient.LLen(pendingTaskQueueName).Result()
	assert.Nil(t, err)
	assert.Empty(t, pendingTaskQueueDepth)

	// Put a task with no execute time on the worker's watched task queue
	task := async.NewTask("foo", nil)
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
	w := newWorker(redisClient, uuid.NewV4().String()).(*worker)

	pendingTaskQueueName := getDisposableQueueName()
	watchedTaskQueueName := getWatchedTaskQueueName(w.id)

	// Assert that the pending task queue is empty
	pendingTaskQueueDepth, err := redisClient.LLen(pendingTaskQueueName).Result()
	assert.Nil(t, err)
	assert.Empty(t, pendingTaskQueueDepth)

	// Put a lapsed task on the worker's watched task queue
	task := async.NewDelayedTask("foo", nil, time.Second*-1)
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
	w := newWorker(redisClient, uuid.NewV4().String()).(*worker)

	pendingTaskQueueName := getDisposableQueueName()
	watchedTaskQueueName := getWatchedTaskQueueName(w.id)

	// Assert that the pending task queue is empty
	pendingTaskQueueDepth, err := redisClient.LLen(pendingTaskQueueName).Result()
	assert.Nil(t, err)
	assert.Empty(t, pendingTaskQueueDepth)

	// Put task with a future execute time on the worker's watched tasks queue
	task := async.NewDelayedTask("foo", nil, time.Second*5)
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

package redis

import (
	"context"
	"testing"
	"time"

	"open-service-broker-azure/pkg/async"
	"github.com/stretchr/testify/assert"
)

func TestDefaultWatchDeferredTasks(t *testing.T) {
	e := getTestEngine()

	pendingTaskQueueName := getDisposableQueueName()
	watchedTaskQueueName := getWatchedTaskQueueName(e.workerID)

	// Define some tasks
	invalidTaskJSON := []byte("bogus")
	taskWithNoExecuteTime := async.NewTask("foo", nil)
	taskWithNoExecuteTimeJSON, err := taskWithNoExecuteTime.ToJSON()
	assert.Nil(t, err)
	validTask := async.NewDelayedTask("foo", nil, time.Second)
	validTaskJSON, err := validTask.ToJSON()
	assert.Nil(t, err)
	tasks := [][]byte{
		invalidTaskJSON,
		taskWithNoExecuteTimeJSON,
		validTaskJSON,
	}

	// Put all the tasks on the worker's watched task queue
	for _, task := range tasks {
		err = e.redisClient.LPush(watchedTaskQueueName, task).Err()
		assert.Nil(t, err)
	}

	// Assert that the pending task queue is empty
	pendingTaskQueueDepth, err := e.redisClient.LLen(pendingTaskQueueName).Result()
	assert.Nil(t, err)
	assert.Empty(t, pendingTaskQueueDepth)

	// Assert that worker's watched task queue has precisely len(tasks) tasks
	watchedTaskQueueDepth, err := e.redisClient.LLen(watchedTaskQueueName).Result()
	assert.Nil(t, err)
	assert.Equal(t, int64(len(tasks)), watchedTaskQueueDepth)

	// Under nominal conditions, defaultWatchDeferredTasks blocks until the
	// context it is passed is canceled. Use a context that will cancel itself
	// after 2 seconds to make defaultWatchDeferredTasks STOP working so we can
	// then examine what it accomplished. (It's set to 2 seconds because the one
	// valid defered task we're executing has a 1 second delay.)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// Put all the tasks on the defaultWatchDeferredTasks function's input channel
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
	go e.defaultWatchDeferredTasks(
		ctx,
		inputCh,
		pendingTaskQueueName,
		errCh,
	)

	select {
	case <-errCh:
		assert.Fail(t, "should not have received any error, but did")
	case <-ctx.Done():
	}

	// Assert that the pending task queue has precisely one task
	pendingTaskQueueDepth, err = e.redisClient.LLen(pendingTaskQueueName).Result()
	assert.Nil(t, err)
	assert.Equal(t, int64(1), pendingTaskQueueDepth)

	// Assert that the worker's watched task queue is empty-- in all cases, the
	// tasks should have been removed from this queue
	watchedTaskQueueDepth, err = e.redisClient.LLen(watchedTaskQueueName).Result()
	assert.Nil(t, err)
	assert.Empty(t, watchedTaskQueueDepth)
}

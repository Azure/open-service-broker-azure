package redis

import (
	"context"
	"errors"
	"testing"
	"time"

	"open-service-broker-azure/pkg/async"
	"github.com/stretchr/testify/assert"
)

func TestDefaultExecuteTasks(t *testing.T) {
	e := getTestEngine()

	pendingTaskQueueName := getDisposableQueueName()
	deferredTaskQueueName := getDisposableQueueName()
	activeTaskQueueName := getActiveTaskQueueName(e.workerID)

	// Register some jobs with the worker
	var badJobCallCount int
	err := e.RegisterJob(
		"badJob",
		func(_ context.Context, _ async.Task) ([]async.Task, error) {
			badJobCallCount++
			return nil, errors.New("a deliberate error")
		},
	)
	assert.Nil(t, err)
	var goodJobCallCount int
	err = e.RegisterJob(
		"goodJob",
		func(_ context.Context, _ async.Task) ([]async.Task, error) {
			goodJobCallCount++
			return []async.Task{
				async.NewTask("followUpJob", nil),
				async.NewDelayedTask("followUpJob", nil, time.Minute),
			}, nil
		},
	)
	assert.Nil(t, err)

	// Define some tasks
	invalidTaskJSON := []byte("bogus")
	unregisteredTask := async.NewTask("nonExistingJob", map[string]string{})
	unregisteredTaskJSON, err := unregisteredTask.ToJSON()
	assert.Nil(t, err)
	badTask := async.NewTask("badJob", map[string]string{})
	badTaskJSON, err := badTask.ToJSON()
	assert.Nil(t, err)
	goodTask := async.NewTask("goodJob", map[string]string{})
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
		err := e.redisClient.LPush(activeTaskQueueName, task).Err()
		assert.Nil(t, err)
	}

	// Assert that the pending task queue is empty
	pendingTaskQueueDepth, err := e.redisClient.LLen(pendingTaskQueueName).Result()
	assert.Nil(t, err)
	assert.Empty(t, pendingTaskQueueDepth)

	// Assert that the deferred task queue is empty
	deferredTaskQueueDepth, err :=
		e.redisClient.LLen(deferredTaskQueueName).Result()
	assert.Nil(t, err)
	assert.Empty(t, deferredTaskQueueDepth)

	// Assert that worker's active task queue has precisely len(tasks) tasks
	activeTaskQueueDepth, err := e.redisClient.LLen(activeTaskQueueName).Result()
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
	go e.defaultExecuteTasks(
		ctx,
		inputCh,
		pendingTaskQueueName,
		deferredTaskQueueName,
		errCh,
	)

	select {
	case <-errCh:
		assert.Fail(t, "should not have received any error, but did")
	case <-ctx.Done():
	}

	// Assert that the pending task queue has precisely two tasks-- the
	// unprocessable task, which should have been returned to it, and a follow-up
	// task
	pendingTaskQueueDepth, err = e.redisClient.LLen(pendingTaskQueueName).Result()
	assert.Nil(t, err)
	assert.Equal(t, int64(2), pendingTaskQueueDepth)

	// Assert that the deferred task queue has precisely one task-- a follow-up
	// task
	deferredTaskQueueDepth, err =
		e.redisClient.LLen(deferredTaskQueueName).Result()
	assert.Nil(t, err)
	assert.Equal(t, int64(1), deferredTaskQueueDepth)

	// Assert that the worker's active task queue is empty-- in all cases, the
	// tasks should have been removed from this queue
	activeTaskQueueDepth, err = e.redisClient.LLen(activeTaskQueueName).Result()
	assert.Nil(t, err)
	assert.Empty(t, activeTaskQueueDepth)

	// Assert that the indicated jobs were invoked the appropriate number of times
	assert.Equal(t, 1, badJobCallCount)
	assert.Equal(t, 1, goodJobCallCount)
}

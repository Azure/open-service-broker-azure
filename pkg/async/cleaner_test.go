package async

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCleanerCleanBlocksUntilCleanInternalErrors(t *testing.T) {
	c := newCleaner(redisClient).(*cleaner)
	c.clean = func(string, string, string) error {
		return errSome
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	err := c.Clean(ctx)
	assert.Equal(t, &errCleaning{err: errSome}, err)
}

func TestCleanerCleanBlocksUntilContextCanceled(t *testing.T) {
	c := newCleaner(redisClient).(*cleaner)
	c.clean = func(string, string, string) error {
		return nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	err := c.Clean(ctx)
	assert.Equal(t, ctx.Err(), err)
}

func TestCleanerCleanInternalCleansDeadWorkers(t *testing.T) {
	activeQueueName := getDisposableQueueName()
	delayedQueueName := getDisposableQueueName()
	workerSetName := getDisposableWorkerSetName()
	const expectedCount = 5
	for range [expectedCount]struct{}{} {
		intCmd := redisClient.SAdd(workerSetName, getDisposableWorkerID())
		assert.Nil(t, intCmd.Err())
	}
	c := newCleaner(redisClient).(*cleaner)
	var cleanWorkerCallCount int
	c.cleanWorker = func(string, string, string) error {
		cleanWorkerCallCount++
		return nil
	}
	err := c.clean(workerSetName, activeQueueName, delayedQueueName)
	assert.Nil(t, err)
	assert.Equal(t, expectedCount, cleanWorkerCallCount)
}

func TestCleanerCleanInternalDoesNotCleanLiveWorkers(t *testing.T) {
	mainActiveWorkQueueName := getDisposableQueueName()
	mainDelayedWorkQueueName := getDisposableQueueName()
	workerSetName := getDisposableWorkerSetName()
	for range [5]struct{}{} {
		workerID := getDisposableWorkerID()
		intCmd := redisClient.SAdd(workerSetName, workerID)
		assert.Nil(t, intCmd.Err())
		statusCmd := redisClient.Set(getHeartbeatKey(workerID), aliveIndicator, 0)
		assert.Nil(t, statusCmd.Err())
	}
	c := newCleaner(redisClient).(*cleaner)
	var cleanWorkerCallCount int
	c.cleanWorker = func(string, string, string) error {
		cleanWorkerCallCount++
		return nil
	}
	err := c.clean(
		workerSetName,
		mainActiveWorkQueueName,
		mainDelayedWorkQueueName,
	)
	assert.Nil(t, err)
	assert.Equal(t, 0, cleanWorkerCallCount)
}

func TestCleanerCleanWorker(t *testing.T) {
	mainActiveWorkQueueName := getDisposableQueueName()
	mainDelayedWorkQueueName := getDisposableQueueName()

	// Assert that the main active work queue starts out empty
	intCmd := redisClient.LLen(mainActiveWorkQueueName)
	assert.Nil(t, intCmd.Err())
	mainActiveWorkQueueDepth, err := intCmd.Result()
	assert.Nil(t, err)
	assert.Empty(t, mainActiveWorkQueueDepth)
	// Assert that the main delayed work queue also starts out empty
	intCmd = redisClient.LLen(mainDelayedWorkQueueName)
	assert.Nil(t, intCmd.Err())
	mainDelayedWorkQueueDepth, err := intCmd.Result()
	assert.Nil(t, err)
	assert.Empty(t, mainDelayedWorkQueueDepth)

	workerID := getDisposableWorkerID()
	workerActiveWorkQueueName := getWorkerActiveQueueName(workerID)
	workerDelayedWorkQueueName := getWorkerDelayedQueueName(workerID)

	const taskCount int64 = 5
	for range [taskCount]struct{}{} {
		// Put some dummy tasks onto the worker's active work queue
		intCmd = redisClient.LPush(workerActiveWorkQueueName, "foo")
		assert.Nil(t, intCmd.Err())
		// Also put some dummy tasks onto the worker's delayed work queue
		intCmd = redisClient.LPush(workerDelayedWorkQueueName, "foo")
		assert.Nil(t, intCmd.Err())
	}

	// Assert that the worker's active work queue is taskCount deep
	intCmd = redisClient.LLen(workerActiveWorkQueueName)
	assert.Nil(t, intCmd.Err())
	workerActiveWorkQueueDepth, err := intCmd.Result()
	assert.Nil(t, err)
	assert.Equal(t, taskCount, workerActiveWorkQueueDepth)
	// Assert that the worker's delayed work queue also is taskCount deep
	intCmd = redisClient.LLen(workerDelayedWorkQueueName)
	assert.Nil(t, intCmd.Err())
	workerDelayedWorkQueueDepth, err := intCmd.Result()
	assert.Nil(t, err)
	assert.Equal(t, taskCount, workerDelayedWorkQueueDepth)

	c := newCleaner(redisClient).(*cleaner)
	err = c.cleanWorker(
		workerID,
		mainActiveWorkQueueName,
		mainDelayedWorkQueueName,
	)
	assert.Nil(t, err)

	// Assert that the main active work queue is now taskCount deep
	intCmd = redisClient.LLen(mainActiveWorkQueueName)
	assert.Nil(t, intCmd.Err())
	mainActiveWorkQueueDepth, err = intCmd.Result()
	assert.Nil(t, err)
	assert.Equal(t, taskCount, mainActiveWorkQueueDepth)
	// Assert that the main delayed work queue also is now taskCount deep
	intCmd = redisClient.LLen(mainDelayedWorkQueueName)
	assert.Nil(t, intCmd.Err())
	mainDelayedWorkQueueDepth, err = intCmd.Result()
	assert.Nil(t, err)
	assert.Equal(t, taskCount, mainDelayedWorkQueueDepth)

	// Assert that the worker's active work queue is now empty
	intCmd = redisClient.LLen(workerActiveWorkQueueName)
	assert.Nil(t, intCmd.Err())
	workerActiveWorkQueueDepth, err = intCmd.Result()
	assert.Nil(t, err)
	assert.Empty(t, workerActiveWorkQueueDepth)
	// Assert that the worker's delayed work queue also is now empty
	intCmd = redisClient.LLen(workerDelayedWorkQueueName)
	assert.Nil(t, intCmd.Err())
	workerDelayedWorkQueueDepth, err = intCmd.Result()
	assert.Nil(t, err)
	assert.Empty(t, workerDelayedWorkQueueDepth)
}

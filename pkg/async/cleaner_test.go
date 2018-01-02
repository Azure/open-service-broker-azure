package async

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCleanerCleanBlocksUntilCleanInternalErrors(t *testing.T) {
	c := newCleaner(redisClient).(*cleaner)
	c.clean = func(string, string) error {
		return errSome
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	err := c.Clean(ctx)
	assert.Equal(t, &errCleaning{err: errSome}, err)
}

func TestCleanerCleanBlocksUntilContextCanceled(t *testing.T) {
	c := newCleaner(redisClient).(*cleaner)
	c.clean = func(string, string) error {
		return nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	err := c.Clean(ctx)
	assert.Equal(t, ctx.Err(), err)
}

func TestCleanerCleanInternalCleansDeadWorkers(t *testing.T) {
	queueName := getDisposableQueueName()
	workerSetName := getDisposableWorkerSetName()
	const expectedCount = 5
	for range [expectedCount]struct{}{} {
		intCmd := redisClient.SAdd(workerSetName, getDisposableWorkerID())
		assert.Nil(t, intCmd.Err())
	}
	c := newCleaner(redisClient).(*cleaner)
	var cleanWorkerCallCount int
	c.cleanWorker = func(string, string) error {
		cleanWorkerCallCount++
		return nil
	}
	err := c.clean(workerSetName, queueName)
	assert.Nil(t, err)
	assert.Equal(t, expectedCount, cleanWorkerCallCount)
}

func TestCleanerCleanInternalDoesNotCleanLiveWorkers(t *testing.T) {
	mainQueueName := getDisposableQueueName()
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
	c.cleanWorker = func(string, string) error {
		cleanWorkerCallCount++
		return nil
	}
	err := c.clean(workerSetName, mainQueueName)
	assert.Nil(t, err)
	assert.Equal(t, 0, cleanWorkerCallCount)
}

func TestCleanerCleanWorker(t *testing.T) {
	mainQueueName := getDisposableQueueName()
	workerID := getDisposableWorkerID()
	workerQueueName := getWorkerQueueName(workerID)
	const taskCount = 5
	for range [taskCount]struct{}{} {
		intCmd := redisClient.LPush(workerQueueName, "foo")
		assert.Nil(t, intCmd.Err())
	}
	c := newCleaner(redisClient).(*cleaner)
	err := c.cleanWorker(workerID, mainQueueName)
	assert.Nil(t, err)
	intCmd := redisClient.LLen(mainQueueName)
	assert.Nil(t, intCmd.Err())
	mainQueueDepth, err := intCmd.Result()
	assert.Nil(t, err)
	assert.Equal(t, int64(taskCount), mainQueueDepth)
	intCmd = redisClient.LLen(workerQueueName)
	assert.Nil(t, intCmd.Err())
	workerQueueDepth, err := intCmd.Result()
	assert.Nil(t, err)
	assert.Empty(t, workerQueueDepth)
}

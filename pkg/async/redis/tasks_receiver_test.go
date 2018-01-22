package redis

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDefaultReceiveTasks(t *testing.T) {
	e := getTestEngine()

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
	go e.defaultReceiveTasks(
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

package redis

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestBeatError(t *testing.T) {
	h := newHeart(getDisposableWorkerID(), time.Second, redisClient).(*heart)

	// Override the default beat function to just returns an error
	h.beat = func() error {
		return errSome
	}

	err := h.Beat()

	// Assert that the error returned from the Beat function wraps the error that
	// the fake beat function generated
	assert.Equal(t, &errHeartbeat{workerID: h.workerID, err: errSome}, err)
}

func TestHeartRunBlocksUntilBeatErrors(t *testing.T) {
	h := newHeart(getDisposableWorkerID(), time.Second, redisClient).(*heart)

	h.beat = func() error {
		return errSome
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Call Run in a goroutine. If it never unblocks, as we hope it does, we don't
	// want the test to stall.
	errCh := make(chan error)
	go func() {
		errCh <- h.Run(ctx)
	}()

	// Assert that the error received from the Run function wraps the error that
	// the overridden watchDeferredTask function generated
	select {
	case err := <-errCh:
		assert.Equal(t, &errHeartbeat{workerID: h.workerID, err: errSome}, err)
	case <-time.After(time.Second):
		assert.Fail(t, "an error should have been received, but wasn't")
	}
}

func TestHeartRunRespondsToContextCanceled(t *testing.T) {
	h := newHeart(getDisposableWorkerID(), time.Second, redisClient).(*heart)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Call Run in a goroutine. If it never unblocks, as we hope it does, we don't
	// want the test to stall.
	errCh := make(chan error)
	go func() {
		errCh <- h.Run(ctx)
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

// TestDefaultBeat tests the happy path for sending a single heartbeat. The
// expected result is that the heartnbeat is visible, with a TTL, in Redis.
func TestDefaultBeat(t *testing.T) {
	h := newHeart(getDisposableWorkerID(), time.Second, redisClient).(*heart)

	err := h.defaultBeat()
	assert.Nil(t, err)

	// Assert that the heartbeat is visible, with a TTL, in Redis.
	str, err := redisClient.Get(getHeartbeatKey(h.workerID)).Result()
	assert.Nil(t, err)
	assert.Equal(t, aliveIndicator, str)
	ttl, err := redisClient.TTL(getHeartbeatKey(h.workerID)).Result()
	assert.Nil(t, err)
	assert.True(t, ttl > 0)
}

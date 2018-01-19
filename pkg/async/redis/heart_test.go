package redis

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDefaultRunHeartBlocksUntilBeatErrors(t *testing.T) {
	w := newWorker(redisClient).(*worker)

	// Override default heartbeat function so it just returns an error
	w.heartbeat = func() error {
		return errSome
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Call defaultRunHeart in a goroutine. If it never unblocks, as we hope it
	// does, we don't want the test to stall.
	errCh := make(chan error)
	go func() {
		errCh <- w.defaultRunHeart(ctx)
	}()

	// Assert that the error received from the defaultRunHeart function is the
	// error that the overridden heartbeat function generated
	select {
	case err := <-errCh:
		assert.Equal(t, errSome, err)
	case <-time.After(time.Second):
		assert.Fail(t, "an error should have been received, but wasn't")
	}
}

func TestDefaultRunHeartRespondsToContextCanceled(t *testing.T) {
	w := newWorker(redisClient).(*worker)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Call defaultRunHeart in a goroutine. If it never unblocks, as we hope it
	// does, we don't want the test to stall.
	errCh := make(chan error)
	go func() {
		errCh <- w.defaultRunHeart(ctx)
	}()

	cancel()

	// Assert that the error returned from defaultRunHeart indicates that the
	// context was canceled
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

// TestDefaultHeartbeat tests the happy path for sending a single heartbeat. The
// expected result is that the heartbeat is visible, with a TTL, in Redis.
func TestDefaultHeartbeat(t *testing.T) {
	w := newWorker(redisClient).(*worker)

	err := w.defaultHeartbeat()
	assert.Nil(t, err)

	// Assert that the heartbeat is visible, with a TTL, in Redis.
	str, err := redisClient.Get(getHeartbeatKey(w.id)).Result()
	assert.Nil(t, err)
	assert.Equal(t, aliveIndicator, str)
	ttl, err := redisClient.TTL(getHeartbeatKey(w.id)).Result()
	assert.Nil(t, err)
	assert.True(t, ttl > 0)
}

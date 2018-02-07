package redis

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDefaultRunHeartBlocksUntilBeatErrors(t *testing.T) {
	e := NewEngine(redisClient).(*engine)

	// Override default heartbeat function so it just returns an error
	e.heartbeat = func(time.Duration) error {
		return errSome
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Call defaultRunHeart in a goroutine. If it never unblocks, as we hope it
	// does, we don't want the test to stall.
	errCh := make(chan error)
	go func() {
		errCh <- e.defaultRunHeart(ctx, time.Second)
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

func TestDefaultRunHeartRespondsToCanceledContext(t *testing.T) {
	e := NewEngine(redisClient).(*engine)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Call defaultRunHeart in a goroutine. If it never unblocks, as we hope it
	// does, we don't want the test to stall.
	errCh := make(chan error)
	go func() {
		errCh <- e.defaultRunHeart(ctx, time.Second)
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

func TestDefaultHeartbeat(t *testing.T) {
	e := NewEngine(redisClient).(*engine)

	err := e.defaultHeartbeat(time.Second)
	assert.Nil(t, err)

	// Assert that the heartbeat is visible, with a TTL, in Redis.
	str, err := redisClient.Get(getHeartbeatKey(e.workerID)).Result()
	assert.Nil(t, err)
	assert.Equal(t, aliveIndicator, str)
	ttl, err := redisClient.TTL(getHeartbeatKey(e.workerID)).Result()
	assert.Nil(t, err)
	assert.True(t, ttl > 0)
}

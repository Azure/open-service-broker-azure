package async

import (
	"context"
	"testing"

	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

func TestHeartBeatError(t *testing.T) {
	workerID := uuid.NewV4().String()
	h := newHeart(workerID, time.Second, redisClient).(*heart)
	h.beat = func() error {
		return errSome
	}
	err := h.Beat()
	assert.Equal(t, &errHeartbeat{workerID: workerID, err: errSome}, err)
}

func TestHeartBeat(t *testing.T) {
	h := newHeart(uuid.NewV4().String(), time.Second, redisClient).(*heart)
	err := h.Beat()
	assert.Nil(t, err)
	strCmd := redisClient.Get(h.workerID)
	assert.Nil(t, strCmd.Err())
	str, err := strCmd.Result()
	assert.Nil(t, err)
	assert.Equal(t, aliveIndicator, str)
}

func TestHeartStartBlocksUntilBeatErrors(t *testing.T) {
	workerID := uuid.NewV4().String()
	h := newHeart(workerID, time.Second, redisClient).(*heart)
	h.beat = func() error {
		return errSome
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	err := h.Start(ctx)
	assert.Equal(t, &errHeartbeat{workerID: workerID, err: errSome}, err)
}

func TestHeartStartBlocksUntilContextCanceled(t *testing.T) {
	h := newHeart(uuid.NewV4().String(), time.Second, redisClient).(*heart)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	err := h.Start(ctx)
	assert.Equal(t, ctx.Err(), err)
}

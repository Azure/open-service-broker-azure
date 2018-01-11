package async

import (
	"context"
	"testing"
	"time"

	"github.com/Azure/open-service-broker-azure/pkg/async/model"
	"github.com/stretchr/testify/assert"
)

func TestResumerResumeBlocksUntilResumeInternalErrors(t *testing.T) {
	r := newResumer(redisClient).(*resumer)
	r.resume = func(q string) error {
		return errSome
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	err := r.Resume(ctx)
	assert.Equal(t, &errResuming{err: errSome}, err)
}

func TestResumerResumeBlocksUntilContextCanceled(t *testing.T) {
	r := newResumer(redisClient).(*resumer)
	r.resume = func(q string) error {
		return nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	err := r.Resume(ctx)
	assert.Equal(t, ctx.Err(), err)
}

func TestResumerEmptiesWatchedQueues(t *testing.T) {
	identifier := getDisposableQueueName()
	r := newResumer(redisClient).(*resumer)
	r.Watch(identifier)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	assert.Equal(t, len(r.watchedQueues), 1)
	_ = r.Resume(ctx)
	assert.Equal(t, len(r.watchedQueues), 0)
}

func TestResumerMoviesTasks(t *testing.T) {
	identifier := getDisposableQueueName()
	r := newResumer(redisClient).(*resumer)
	task := model.NewTask(
		"waitForParentStep",
		map[string]string{
			"provisionFirstStep": "asdf",
			"instanceID":         "asdf",
		},
	)
	json, _ := task.ToJSON()
	redisClient.LPush(identifier, json)
	intCmd := redisClient.LLen(mainWorkQueueName)
	assert.Nil(t, intCmd.Err())
	mainQueueDepth, _ := intCmd.Result()
	assert.Equal(t, int64(0), mainQueueDepth)
	r.Watch(identifier)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	_ = r.Resume(ctx)
	intCmd = redisClient.LLen(mainWorkQueueName)
	assert.Nil(t, intCmd.Err())
	mainQueueDepth, _ = intCmd.Result()
	assert.Equal(t, int64(1), mainQueueDepth)
}

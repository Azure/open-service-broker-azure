package async

import (
	"context"
	"testing"
	"time"

	fakeAsync "github.com/Azure/azure-service-broker/pkg/async/fake"
	"github.com/Azure/azure-service-broker/pkg/async/model"
	"github.com/stretchr/testify/assert"
)

func TestWorkerGetsUniqueID(t *testing.T) {
	// Create two workers-- make sure their IDs are at least different from one
	// another
	w1 := newWorker(redisClient).(*worker)
	w2 := newWorker(redisClient).(*worker)
	assert.NotEqual(t, w1.id, w2.id)
}

func TestWorkerWorkBlocksUntilHeartErrors(t *testing.T) {
	h := fakeAsync.NewHeart()
	h.RunBehavior = func(context.Context) error {
		return errSome
	}
	w := newWorker(redisClient).(*worker)
	w.heart = h
	receiveAndWorkStopped := false
	w.receiveAndWork = func(ctx context.Context, queueName string) error {
		<-ctx.Done()
		receiveAndWorkStopped = true
		return ctx.Err()
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	err := w.Work(ctx)
	assert.Equal(
		t,
		&errHeartStopped{workerID: w.id, err: errSome},
		err,
	)
	time.Sleep(time.Second)
	assert.True(t, receiveAndWorkStopped)
}

func TestWorkerWorkBlocksUntilHeartReturns(t *testing.T) {
	h := fakeAsync.NewHeart()
	h.RunBehavior = func(context.Context) error {
		return nil
	}
	w := newWorker(redisClient).(*worker)
	w.heart = h
	receiveAndWorkStopped := false
	w.receiveAndWork = func(ctx context.Context, queueName string) error {
		<-ctx.Done()
		receiveAndWorkStopped = true
		return ctx.Err()
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	err := w.Work(ctx)
	assert.Equal(
		t,
		&errHeartStopped{workerID: w.id},
		err,
	)
	time.Sleep(time.Second)
	assert.True(t, receiveAndWorkStopped)
}

func TestWorkerWorkBlocksUntilReceiveAndWorkErrors(t *testing.T) {
	h := fakeAsync.NewHeart()
	heartStopped := false
	h.RunBehavior = func(ctx context.Context) error {
		<-ctx.Done()
		heartStopped = true
		return ctx.Err()
	}
	w := newWorker(redisClient).(*worker)
	w.heart = h
	w.receiveAndWork = func(context.Context, string) error {
		return errSome
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	err := w.Work(ctx)
	assert.Equal(
		t,
		&errReceiveAndWorkStopped{workerID: w.id, err: errSome},
		err,
	)
	time.Sleep(time.Second)
	assert.True(t, heartStopped)
}

func TestWorkerWorkBlocksUntilReceiveAndWorkReturns(t *testing.T) {
	h := fakeAsync.NewHeart()
	heartStopped := false
	h.RunBehavior = func(ctx context.Context) error {
		<-ctx.Done()
		heartStopped = true
		return ctx.Err()
	}
	w := newWorker(redisClient).(*worker)
	w.heart = h
	w.receiveAndWork = func(context.Context, string) error {
		return nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	err := w.Work(ctx)
	assert.Equal(
		t,
		&errReceiveAndWorkStopped{workerID: w.id},
		err,
	)
	time.Sleep(time.Second)
	assert.True(t, heartStopped)
}

func TestWorkerWorkBlocksUntilContextCanceled(t *testing.T) {
	h := fakeAsync.NewHeart()
	heartStopped := false
	h.RunBehavior = func(ctx context.Context) error {
		<-ctx.Done()
		heartStopped = true
		return ctx.Err()
	}
	w := newWorker(redisClient).(*worker)
	w.heart = h
	receiveAndWorkStopped := false
	w.receiveAndWork = func(ctx context.Context, queueName string) error {
		<-ctx.Done()
		receiveAndWorkStopped = true
		return ctx.Err()
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	err := w.Work(ctx)
	assert.Equal(t, ctx.Err(), err)
	time.Sleep(time.Second)
	assert.True(t, heartStopped)
	assert.True(t, receiveAndWorkStopped)
}

func TestReceiveAndWorkCallsWorkOncePerTask(t *testing.T) {
	queueName := getDisposableQueueName()
	const expectedCount = 5
	for range [expectedCount]struct{}{} {
		taskJSON, err := model.NewTask("foo", nil).ToJSON()
		assert.Nil(t, err)
		intCmd := redisClient.LPush(queueName, taskJSON)
		assert.Nil(t, intCmd.Err())
	}
	w := newWorker(redisClient).(*worker)
	var workCount int
	w.work = func(context.Context, model.Task) error {
		workCount++
		return nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	err := w.receiveAndWork(ctx, queueName)
	assert.Equal(t, ctx.Err(), err)
	assert.Equal(t, expectedCount, workCount)
	intCmd := redisClient.LLen(queueName)
	assert.Nil(t, intCmd.Err())
	currentMainQueueDepth, err := intCmd.Result()
	assert.Nil(t, err)
	assert.Empty(t, currentMainQueueDepth)
	intCmd = redisClient.LLen(getWorkerQueueName(w.id))
	assert.Nil(t, intCmd.Err())
	currentWorkerQueueDepth, err := intCmd.Result()
	assert.Nil(t, err)
	assert.Empty(t, currentWorkerQueueDepth)
}

func TestWorkerReceiveAndWorkBlocksEvenAfterInvalidTask(t *testing.T) {
	queueName := getDisposableQueueName()
	intCmd := redisClient.LPush(queueName, "bogus")
	assert.Nil(t, intCmd.Err())
	w := newWorker(redisClient).(*worker)
	workCalled := false
	w.work = func(context.Context, model.Task) error {
		workCalled = true
		return nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	err := w.receiveAndWork(ctx, queueName)
	assert.Equal(t, ctx.Err(), err)
	assert.False(t, workCalled)
	intCmd = redisClient.LLen(queueName)
	assert.Nil(t, intCmd.Err())
	currentMainQueueDepth, err := intCmd.Result()
	assert.Nil(t, err)
	assert.Empty(t, currentMainQueueDepth)
	intCmd = redisClient.LLen(getWorkerQueueName(w.id))
	assert.Nil(t, intCmd.Err())
	currentWorkerQueueDepth, err := intCmd.Result()
	assert.Nil(t, err)
	assert.Empty(t, currentWorkerQueueDepth)
}

func TestWorkerReceiveAndWorkBlocksEvenAfterWorkError(t *testing.T) {
	queueName := getDisposableQueueName()
	taskJSON, err := model.NewTask("foo", nil).ToJSON()
	assert.Nil(t, err)
	intCmd := redisClient.LPush(queueName, taskJSON)
	assert.Nil(t, intCmd.Err())
	w := newWorker(redisClient).(*worker)
	workCalled := false
	w.work = func(context.Context, model.Task) error {
		workCalled = true
		return errSome
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	err = w.receiveAndWork(ctx, queueName)
	assert.Equal(t, ctx.Err(), err)
	assert.True(t, workCalled)
}

func TestWorkerReceiveAndWorkBlocksUntilContextCanceled(t *testing.T) {
	w := newWorker(redisClient).(*worker)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	err := w.receiveAndWork(ctx, getDisposableQueueName())
	assert.Equal(t, ctx.Err(), err)
}

package async

import (
	"context"
	"testing"
	"time"

	fakeAsync "github.com/Azure/azure-service-broker/pkg/async/fake"
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
	var receiveAndWorkStopped bool
	w.receiveAndWork = func(ctx context.Context) error {
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
	var receiveAndWorkStopped bool
	w.receiveAndWork = func(ctx context.Context) error {
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
	var heartStopped bool
	h.RunBehavior = func(ctx context.Context) error {
		<-ctx.Done()
		heartStopped = true
		return ctx.Err()
	}
	w := newWorker(redisClient).(*worker)
	w.heart = h
	w.receiveAndWork = func(context.Context) error {
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
	var heartStopped bool
	h.RunBehavior = func(ctx context.Context) error {
		<-ctx.Done()
		heartStopped = true
		return ctx.Err()
	}
	w := newWorker(redisClient).(*worker)
	w.heart = h
	w.receiveAndWork = func(context.Context) error {
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
	var heartStopped bool
	h.RunBehavior = func(ctx context.Context) error {
		<-ctx.Done()
		heartStopped = true
		return ctx.Err()
	}
	w := newWorker(redisClient).(*worker)
	w.heart = h
	var receiveAndWorkStopped bool
	w.receiveAndWork = func(ctx context.Context) error {
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

func TestWorkerReceiveAndWorkBlocksUntilError(t *testing.T) {
	// TODO: Implement this
}

func TestWorkerReceiveAndWorkBlocksUntilContextCanceled(t *testing.T) {
	w := newWorker(redisClient).(*worker)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	err := w.receiveAndWork(ctx)
	assert.Equal(t, ctx.Err(), err)
}

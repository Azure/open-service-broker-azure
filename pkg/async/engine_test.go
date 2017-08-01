package async

import (
	"context"
	"testing"
	"time"

	fakeAsync "github.com/Azure/azure-service-broker/pkg/async/fake"
	"github.com/stretchr/testify/assert"
)

func TestEngineStartBlocksUntilCleanerErrors(t *testing.T) {
	e := NewEngine(redisClient).(*engine)
	c := fakeAsync.NewCleaner()
	c.RunBehavior = func(context.Context) error {
		return errSome
	}
	e.cleaner = c
	workerStopped := false
	w := fakeAsync.NewWorker()
	w.RunBehavior = func(ctx context.Context) error {
		<-ctx.Done()
		workerStopped = true
		return ctx.Err()
	}
	e.worker = w
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	err := e.Start(ctx)
	assert.Equal(t, &errCleanerStopped{err: errSome}, err)
	time.Sleep(time.Second)
	assert.True(t, workerStopped)
}

func TestEngineStartBlocksUntilCleanerReturns(t *testing.T) {
	e := NewEngine(redisClient).(*engine)
	c := fakeAsync.NewCleaner()
	c.RunBehavior = func(context.Context) error {
		return nil
	}
	e.cleaner = c
	workerStopped := false
	w := fakeAsync.NewWorker()
	w.RunBehavior = func(ctx context.Context) error {
		<-ctx.Done()
		workerStopped = true
		return ctx.Err()
	}
	e.worker = w
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	err := e.Start(ctx)
	assert.Equal(t, &errCleanerStopped{}, err)
	time.Sleep(time.Second)
	assert.True(t, workerStopped)
}

func TestEngineStartBlocksUntilWorkerErrors(t *testing.T) {
	e := NewEngine(redisClient).(*engine)
	cleanerStopped := false
	c := fakeAsync.NewCleaner()
	c.RunBehavior = func(ctx context.Context) error {
		<-ctx.Done()
		cleanerStopped = true
		return ctx.Err()
	}
	e.cleaner = c
	w := fakeAsync.NewWorker()
	w.RunBehavior = func(context.Context) error {
		return errSome
	}
	e.worker = w
	err := e.Start(context.Background())
	assert.Equal(t, &errWorkerStopped{workerID: w.GetID(), err: errSome}, err)
	time.Sleep(time.Second)
	assert.True(t, cleanerStopped)
}

func TestEngineStartBlocksUntilWorkerReturns(t *testing.T) {
	e := NewEngine(redisClient).(*engine)
	cleanerStopped := false
	c := fakeAsync.NewCleaner()
	c.RunBehavior = func(ctx context.Context) error {
		<-ctx.Done()
		cleanerStopped = true
		return ctx.Err()
	}
	e.cleaner = c
	w := fakeAsync.NewWorker()
	w.RunBehavior = func(context.Context) error {
		return nil
	}
	e.worker = w
	err := e.Start(context.Background())
	assert.Equal(t, &errWorkerStopped{workerID: w.GetID()}, err)
	time.Sleep(time.Second)
	assert.True(t, cleanerStopped)
}

func TestEngineStartBlocksUntilContextCanceled(t *testing.T) {
	e := NewEngine(redisClient).(*engine)
	cleanerStopped := false
	c := fakeAsync.NewCleaner()
	c.RunBehavior = func(ctx context.Context) error {
		<-ctx.Done()
		cleanerStopped = true
		return ctx.Err()
	}
	e.cleaner = c
	workerStopped := false
	w := fakeAsync.NewWorker()
	w.RunBehavior = func(ctx context.Context) error {
		<-ctx.Done()
		workerStopped = true
		return ctx.Err()
	}
	e.worker = w
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	err := e.Start(ctx)
	assert.Equal(t, ctx.Err(), err)
	time.Sleep(time.Second)
	assert.True(t, cleanerStopped)
	assert.True(t, workerStopped)
}

package async

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCleanerRunBlocksUntilCleanErrors(t *testing.T) {
	c := newCleaner(redisClient).(*cleaner)

	// Override the default clean function to just return an error
	c.clean = func(context.Context, string, string, string) error {
		return errSome
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Call Run in a goroutine. If it never unblocks, as we hope it does, we don't
	// want the test to stall.
	errCh := make(chan error)
	go func() {
		errCh <- c.Run(ctx)
	}()

	// Assert that the error returned from the Run function is the error that
	// the overridden clean function generated
	select {
	case err := <-errCh:
		assert.Equal(t, &errCleaning{err: errSome}, err)
	case <-time.After(time.Second):
		assert.Fail(
			t,
			"an error error should have been returned, but wasn't",
		)
	}
}

func TestCleanerRunRespondsToCanceledContext(t *testing.T) {
	c := newCleaner(redisClient).(*cleaner)

	// Override the default clean function to be a no-op
	c.clean = func(context.Context, string, string, string) error {
		return nil
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Call Run in a goroutine. If it never unblocks, as we hope it does, we don't
	// want the test to stall.
	errCh := make(chan error)
	go func() {
		errCh <- c.Run(ctx)
	}()

	cancel()

	// Assert that the error returned from Work indicates that the context was
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

func TestDefaultCleanCleansDeadWorkers(t *testing.T) {
	c := newCleaner(redisClient).(*cleaner)

	// Add some workers to the worker set, but do not add any heartbeats for these
	// workers. i.e. They should appear dead.
	const workerCount = 5
	workerSetName := getDisposableWorkerSetName()
	for range [workerCount]struct{}{} {
		err := redisClient.SAdd(workerSetName, getDisposableWorkerID()).Err()
		assert.Nil(t, err)
	}

	// Override the default cleanActiveTaskQueue function to just count how many
	// times it is invoked
	var cleanActiveTaskQueueCallCount int
	c.cleanActiveTaskQueue = func(context.Context, string, string, string) error {
		cleanActiveTaskQueueCallCount++
		return nil
	}

	// Override the default cleanWatchedTaskQueue function to just count how many
	// times it is invoked
	var cleanWatchedTaskQueueCallCount int
	c.cleanWatchedTaskQueue = func(
		context.Context,
		string,
		string,
		string,
	) error {
		cleanWatchedTaskQueueCallCount++
		return nil
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := c.defaultClean(
		ctx,
		workerSetName,
		getDisposableQueueName(),
		getDisposableQueueName(),
	)
	assert.Nil(t, err)

	// Assert cleanActiveTaskQueue and cleanWatchedTaskQueue were each invoked
	// once per dead worker
	assert.Equal(t, workerCount, cleanActiveTaskQueueCallCount)
	assert.Equal(t, workerCount, cleanWatchedTaskQueueCallCount)
}

func TestDefaultCleanDoesNotCleanLiveWorkers(t *testing.T) {
	c := newCleaner(redisClient).(*cleaner)

	// Add a worker to the worker set. Also add a heartbeat so this worker appears
	// to be alive.
	workerSetName := getDisposableWorkerSetName()
	workerID := getDisposableWorkerID()
	err := redisClient.SAdd(workerSetName, workerID).Err()
	assert.Nil(t, err)
	err = redisClient.Set(getHeartbeatKey(workerID), aliveIndicator, 0).Err()
	assert.Nil(t, err)

	// Override the default cleanActiveTaskQueue function to just count how many
	// times it is invoked
	var cleanActiveTaskQueueCallCount int
	c.cleanActiveTaskQueue = func(context.Context, string, string, string) error {
		cleanActiveTaskQueueCallCount++
		return nil
	}

	// Override the default cleanWatchedTaskQueue function to just count how many
	// times it is invoked
	var cleanWatchedTaskQueueCallCount int
	c.cleanWatchedTaskQueue = func(
		context.Context,
		string,
		string,
		string,
	) error {
		cleanWatchedTaskQueueCallCount++
		return nil
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err = c.defaultClean(
		ctx,
		workerSetName,
		getDisposableQueueName(),
		getDisposableQueueName(),
	)
	assert.Nil(t, err)

	// Assert neither cleanActiveTaskQueue and cleanWatchedTaskQueue were ever
	// invoked
	assert.Equal(t, 0, cleanActiveTaskQueueCallCount)
	assert.Equal(t, 0, cleanWatchedTaskQueueCallCount)
}

func TestDefaultCleanRespondsToContextCanceled(t *testing.T) {
	c := newCleaner(redisClient).(*cleaner)

	// Add one worker to the worker set. Do not add a heartbeat. This should
	// guarantee that some cleanup must take place. This is imporant because
	// we'll override the default cleanActiveTaskQueue function to block us for a
	// while, giving us the opportunity to test that defaultClean responds to
	// context cancelation.
	workerSetName := getDisposableWorkerSetName()
	workerID := getDisposableWorkerID()
	err := redisClient.SAdd(workerSetName, workerID).Err()
	assert.Nil(t, err)

	// Override the default cleanActiveTaskQueue function to block until the
	// context it is passed is canceled.
	c.cleanActiveTaskQueue = func(
		ctx context.Context,
		_ string,
		_ string,
		_ string,
	) error {
		<-ctx.Done()
		return ctx.Err()
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// defaultClean doesn't normally have a possibility of blocking indefinitely,
	// but we've overridden the default cleanActiveTaskQueue function that
	// defaultClean calls. We've done this so we can test our ability to
	// short-circuit defaultClean with a canceled context, but in doing so, we've
	// created the possibility for defaultClean to block indefinitely if it does
	// not responds to the canceled context as we hope it does. So, here, we
	// invoke defaultClean in a goroutine so that, in the worst case, the test
	// won't stall.
	errCh := make(chan error)
	go func() {
		errCh <- c.defaultClean(
			ctx,
			workerSetName,
			getDisposableQueueName(),
			getDisposableQueueName(),
		)
	}()

	cancel()

	// Assert that the error returned from defaultClean indicates that the context
	// was canceled
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

func TestDefaultCleanWorkerQueue(t *testing.T) {
	c := newCleaner(redisClient).(*cleaner)

	sourceQueueName := getDisposableQueueName()
	destinationQueueName := getDisposableQueueName()

	const taskCount int64 = 5
	for range [taskCount]struct{}{} {
		// Put some dummy tasks onto the source queue
		err := redisClient.LPush(sourceQueueName, "foo").Err()
		assert.Nil(t, err)
	}

	// Assert that the source queue is precisely taskCount deep
	sourceQueueDepth, err := redisClient.LLen(sourceQueueName).Result()
	assert.Nil(t, err)
	assert.Equal(t, taskCount, sourceQueueDepth)

	// Assert that the destination queue starts out empty
	destinationQueueDepth, err := redisClient.LLen(destinationQueueName).Result()
	assert.Nil(t, err)
	assert.Empty(t, destinationQueueDepth)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err = c.defaultCleanWorkerQueue(
		ctx,
		getDisposableWorkerID(),
		sourceQueueName,
		destinationQueueName,
	)
	assert.Nil(t, err)

	// Assert that the source queue has been drained
	sourceQueueDepth, err = redisClient.LLen(sourceQueueName).Result()
	assert.Nil(t, err)
	assert.Empty(t, sourceQueueDepth)

	// Assert that the destination queue now has precisely taskCount tasks
	destinationQueueDepth, err = redisClient.LLen(destinationQueueName).Result()
	assert.Nil(t, err)
	assert.Equal(t, taskCount, destinationQueueDepth)
}

func TestDefaultCleanWorkerQueueRespondsToContextCanceled(t *testing.T) {
	c := newCleaner(redisClient).(*cleaner)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Cancel before we even call defaultCleanWorkerQueue. In this case, its
	// the only way we can guarantee the function won't return before we have
	// a chance to cancel the context.
	cancel()

	err := c.defaultCleanWorkerQueue(
		ctx,
		getDisposableWorkerID(),
		getDisposableQueueName(),
		getDisposableQueueName(),
	)

	// Assert that the error returned indicates that the context was canceled
	assert.Equal(t, ctx.Err(), err)
}

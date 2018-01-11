package async

import (
	"context"
	"fmt"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/go-redis/redis"
)

type cleanFn func(
	ctx context.Context,
	workerSetName string,
	pendingTaskQueueName string,
	deferredTaskQueueName string,
) error

type cleanWorkerQueueFn func(
	ctx context.Context,
	workerID string,
	pendingTaskQueueName string,
	deferredTaskQueueName string,
) error

// Cleaner is an interface to be implemented by components that re-queue tasks
// assigned to dead workers
type Cleaner interface {
	// Run causes the cleaner to clean up after dead worker. It blocks until a
	// fatal error is encountered or the context passed to it has been canceled.
	// Run always returns a non-nil error.
	Run(context.Context) error
}

// cleaner is a Redis-based implementation of the Cleaner interface
type cleaner struct {
	redisClient *redis.Client
	// This allows tests to inject an alternative implementation of this function
	clean cleanFn
	// This allows tests to inject an alternative implementation of this function
	cleanActiveTaskQueue cleanWorkerQueueFn
	// This allows tests to inject an alternative implementation of this function
	cleanWatchedTaskQueue cleanWorkerQueueFn
}

func newCleaner(redisClient *redis.Client) Cleaner {
	c := &cleaner{
		redisClient: redisClient,
	}
	c.clean = c.defaultClean
	c.cleanActiveTaskQueue = c.defaultCleanWorkerQueue
	c.cleanWatchedTaskQueue = c.defaultCleanWorkerQueue
	return c
}

// Run causes the cleaner to clean up after dead worker. It blocks until a fatal
// error is encountered or the context passed to it has been canceled. Run
// always returns a non-nil error.
func (c *cleaner) Run(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	ticker := time.NewTicker(time.Second * 10)
	defer ticker.Stop()
	for {
		if err := c.clean(
			ctx,
			workerSetName,
			pendingTaskQueueName,
			deferredTaskQueueName,
		); err != nil {
			return &errCleaning{err: err}
		}
		select {
		case <-ticker.C:
		case <-ctx.Done():
			log.Debug("context canceled; async worker cleaner shutting down")
			return ctx.Err()
		}
	}
}

func (c *cleaner) defaultClean(
	ctx context.Context,
	workerSetName string,
	pendingTaskQueueName string,
	deferredTaskQueueName string,
) error {
	workerIDs, err := c.redisClient.SMembers(workerSetName).Result()
	if err == redis.Nil {
		return nil
	}
	if err != nil {
		return fmt.Errorf("error retrieving workers: %s", err)
	}
	for _, workerID := range workerIDs {
		err := c.redisClient.Get(getHeartbeatKey(workerID)).Err()
		if err == nil {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
				continue
			}
		}
		if err != redis.Nil {
			return fmt.Errorf(
				`error checking health of worker: "%s": %s`,
				workerID,
				err,
			)
		}
		// If we get to here, we have a dead worker on our hands
		if err := c.cleanActiveTaskQueue(
			ctx,
			workerID,
			getActiveTaskQueueName(workerID),
			pendingTaskQueueName,
		); err != nil {
			return err
		}
		if err := c.cleanWatchedTaskQueue(
			ctx,
			workerID,
			getWatchedTaskQueueName(workerID),
			deferredTaskQueueName,
		); err != nil {
			return err
		}
		err = c.redisClient.SRem(workerSetName, workerID).Err()
		if err != nil && err != redis.Nil {
			return fmt.Errorf(
				`error removing dead worker "%s" from worker set: %s`,
				workerID,
				err,
			)
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
	}
	return nil
}

func (c *cleaner) defaultCleanWorkerQueue(
	ctx context.Context,
	workerID string,
	sourceQueueName string,
	destinationQueueName string,
) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		err := c.redisClient.RPopLPush(sourceQueueName, destinationQueueName).Err()
		if err == redis.Nil {
			return nil
		}
		if err != nil {
			return fmt.Errorf(
				`error cleaning up after dead worker "%s" queue "%s": %s`,
				workerID,
				sourceQueueName,
				err,
			)
		}
	}
}

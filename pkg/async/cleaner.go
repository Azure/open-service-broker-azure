package async

import (
	"context"
	"fmt"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/go-redis/redis"
)

type cleanFunction func(workerSetName, mainWorkQueueName string) error

type cleanWorkerFunction func(workerID, mainWorkQueueName string) error

// Cleaner is an interface to be implemented by components that re-queue work
// assigned to dead workers
type Cleaner interface {
	Clean(context.Context) error
}

// cleaner is a Redis-based implementation of the Cleaner interface
type cleaner struct {
	redisClient *redis.Client
	// This allows tests to inject an alternative implementation of this function
	clean cleanFunction
	// This allows tests to inject an alternative implementation of this function
	cleanWorker cleanWorkerFunction
}

func newCleaner(redisClient *redis.Client) Cleaner {
	c := &cleaner{
		redisClient: redisClient,
	}
	c.clean = c.defaultClean
	c.cleanWorker = c.defaultCleanWorker
	return c
}

func (c *cleaner) Clean(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	ticker := time.NewTicker(time.Second * 10)
	defer ticker.Stop()
	for {
		if err := c.clean("workers", mainWorkQueueName); err != nil {
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

func (c *cleaner) defaultClean(workerSetName, mainWorkQueueName string) error {
	strsCmd := c.redisClient.SMembers(workerSetName)
	if strsCmd.Err() == nil {
		workerIDs, err := strsCmd.Result()
		if err != nil {
			return fmt.Errorf("error retrieving workers: %s", err)
		}
		for _, workerID := range workerIDs {
			strCmd := c.redisClient.Get(getHeartbeatKey(workerID))
			if strCmd.Err() == nil {
				continue
			}
			if strCmd.Err() != redis.Nil {
				return fmt.Errorf(
					`error checking health of worker: "%s": %s`,
					workerID,
					strCmd.Err(),
				)
			}
			// If we get to here, we have a dead worker on our hands
			if err := c.cleanWorker(workerID, mainWorkQueueName); err != nil {
				return fmt.Errorf(
					`error cleaning up after dead worker "%s": %s`,
					workerID,
					err,
				)
			}
			intCmd := c.redisClient.SRem("workers", workerID)
			if intCmd.Err() != nil && intCmd.Err() != redis.Nil {
				return fmt.Errorf(
					`error removing dead worker "%s" from worker set: %s`,
					workerID,
					intCmd.Err(),
				)
			}
		}
	} else if strsCmd.Err() != redis.Nil {
		return fmt.Errorf("error retrieving workers: %s", strsCmd.Err())
	}
	return nil
}

func (c *cleaner) defaultCleanWorker(workerID, mainWorkQueueName string) error {
	for {
		strCmd := c.redisClient.RPopLPush(
			getWorkerQueueName(workerID),
			mainWorkQueueName,
		)
		if strCmd.Err() == redis.Nil {
			return nil
		}
		if strCmd.Err() != nil {
			return fmt.Errorf(
				`error cleaning up after dead worker "%s": %s`,
				workerID,
				strCmd.Err(),
			)
		}
	}
}

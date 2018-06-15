package redis

import (
	"context"
	"fmt"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/go-redis/redis"
)

const cleaningInterval = time.Second * 30

type cleanFn func(
	ctx context.Context,
	workerSetName string,
	pendingTaskQueueName string,
	deferredTaskQueueName string,
	interval time.Duration,
) error

type cleanWorkerQueueFn func(
	ctx context.Context,
	workerID string,
	pendingTaskQueueName string,
	deferredTaskQueueName string,
) error

func (e *engine) defaultClean(
	ctx context.Context,
	workerSetName string,
	pendingTaskQueueName string,
	deferredTaskQueueName string,
	interval time.Duration,
) error {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			workerIDs, err := e.redisClient.SMembers(workerSetName).Result()
			if err == redis.Nil {
				continue
			}
			if err != nil {
				return fmt.Errorf("error retrieving workers: %s", err)
			}
			for _, workerID := range workerIDs {
				err := e.redisClient.Get(getHeartbeatKey(workerID)).Err()
				if err == nil {
					continue
				}
				if err != redis.Nil {
					return fmt.Errorf(
						`error checking health of worker: "%s": %s`,
						workerID,
						err,
					)
				}
				// If we get to here, we have a dead worker on our hands
				if err := e.cleanActiveTaskQueue(
					ctx,
					workerID,
					getActiveTaskQueueName(workerID),
					pendingTaskQueueName,
				); err != nil {
					return err
				}
				if err := e.cleanWatchedTaskQueue(
					ctx,
					workerID,
					getWatchedTaskQueueName(workerID),
					deferredTaskQueueName,
				); err != nil {
					return err
				}
				err = e.redisClient.SRem(workerSetName, workerID).Err()
				if err != nil && err != redis.Nil {
					return fmt.Errorf(
						`error removing dead worker "%s" from worker set: %s`,
						workerID,
						err,
					)
				}
			}
		case <-ctx.Done():
			log.Debug("context canceled; async worker cleaner shutting down")
			return ctx.Err()
		}
	}
}

func (e *engine) defaultCleanWorkerQueue(
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
		err := e.redisClient.RPopLPush(sourceQueueName, destinationQueueName).Err()
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

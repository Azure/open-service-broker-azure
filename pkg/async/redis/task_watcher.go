package redis

import (
	"context"
	"fmt"
	"time"

	log "github.com/Sirupsen/logrus"
)

// watchDeferredTaskFn defines functions used to watch a deferred task
type watchDeferredTaskFn func(
	ctx context.Context,
	taskJSON []byte,
	pendingTaskQueueName string,
	errCh chan error,
)

func (w *worker) defaultWatchDeferredTask(
	ctx context.Context,
	taskJSON []byte,
	pendingTaskQueueName string,
	errCh chan error,
) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	task, err := w.getTaskFromJSON(taskJSON, getWatchedTaskQueueName(w.id))
	if err != nil {
		select {
		case errCh <- err:
		case <-ctx.Done():
		}
		return
	}
	if task == nil {
		return
	}
	executeTime := task.GetExecuteTime()
	if executeTime == nil {
		err := w.redisClient.LRem(getWatchedTaskQueueName(w.id), -1, taskJSON).Err()
		if err != nil {
			select {
			case errCh <- fmt.Errorf(
				`error removing task "%s" with no executeTime from queue "%s": %s`,
				task.GetID(),
				getWatchedTaskQueueName(w.id),
				err,
			):
			case <-ctx.Done():
			}
			return
		}
		log.WithFields(log.Fields{
			"task":  task.GetID(),
			"queue": getWatchedTaskQueueName(w.id),
		}).Error("deferred task had no executeTime and was removed from the queue")
		return
	}
	// Note if the duration passed to the timer is 0 or negative, it should go
	// off immediately
	timer := time.NewTimer(time.Until(*executeTime))
	defer timer.Stop()
	select {
	case <-timer.C:
		// Move the task to the pending queue
		pipeline := w.redisClient.TxPipeline()
		pipeline.LPush(pendingTaskQueueName, taskJSON)
		pipeline.LRem(getWatchedTaskQueueName(w.id), -1, taskJSON)
		_, err := pipeline.Exec()
		if err != nil {
			select {
			case errCh <- fmt.Errorf(
				`error moving deferred task "%s" to queue "%s": %s`,
				task.GetID(),
				pendingTaskQueueName,
				err,
			):
			case <-ctx.Done():
			}
		}
	case <-ctx.Done():
	}
}

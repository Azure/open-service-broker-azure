package redis

import (
	"context"
	"fmt"
	"time"

	log "github.com/Sirupsen/logrus"
)

// watchDeferredTasksFn defines functions used to watch a deferred task
type watchDeferredTasksFn func(
	ctx context.Context,
	inputCh chan []byte,
	pendingTaskQueueName string,
	errCh chan error,
)

func (e *engine) defaultWatchDeferredTasks(
	ctx context.Context,
	inputCh chan []byte,
	pendingTaskQueueName string,
	errCh chan error,
) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	for {
		select {
		case taskJSON := <-inputCh:
			task, err := e.getTaskFromJSON(
				taskJSON,
				e.watchedTaskQueueName,
			)
			if err != nil {
				select {
				case errCh <- err:
				case <-ctx.Done():
				}
				return
			}
			if task == nil {
				continue
			}
			executeTime := task.GetExecuteTime()
			if executeTime == nil {
				err := e.redisClient.LRem(
					e.watchedTaskQueueName,
					-1,
					taskJSON,
				).Err()
				if err != nil {
					select {
					case errCh <- fmt.Errorf(
						`error removing task "%s" with no executeTime from queue "%s": %s`,
						task.GetID(),
						e.watchedTaskQueueName,
						err,
					):
					case <-ctx.Done():
					}
					return
				}
				log.WithFields(log.Fields{
					"task":  task.GetID(),
					"queue": e.watchedTaskQueueName,
				}).Error("deferred task had no executeTime and was removed from the queue")
				continue
			}
			// Note if the duration passed to the timer is 0 or negative, it should go
			// off immediately
			timer := time.NewTimer(time.Until(*executeTime))
			defer timer.Stop()
			select {
			case <-timer.C:
				// Move the task to the pending queue
				pipeline := e.redisClient.TxPipeline()
				pipeline.LPush(pendingTaskQueueName, taskJSON)
				pipeline.LRem(e.watchedTaskQueueName, -1, taskJSON)
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
		case <-ctx.Done():
			return
		}
	}
}

package redis

import (
	"context"
	"fmt"

	log "github.com/Sirupsen/logrus"
)

// executeTasksFn defines functions used to execute pending tasks
type executeTasksFn func(
	ctx context.Context,
	inputCh chan []byte,
	pendingTaskQueueName string,
	deferredTaskQueueName string,
	errCh chan error,
)

func (e *engine) defaultExecuteTasks(
	ctx context.Context,
	inputCh chan []byte,
	pendingTaskQueueName string,
	deferredTaskQueueName string,
	errCh chan error,
) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	for {
		select {
		case taskJSON := <-inputCh:
			task, err := e.getTaskFromJSON(
				taskJSON,
				getActiveTaskQueueName(e.workerID),
			)
			if err != nil {
				select {
				case errCh <- err:
					continue
				case <-ctx.Done():
					return
				}
			}
			if task == nil {
				continue
			}
			e.jobsFnsMutex.RLock()
			defer e.jobsFnsMutex.RUnlock()
			jobFn, ok := e.jobsFns[task.GetJobName()]
			if !ok {
				// This worker doesn't know how to process this task. That doesn't mean
				// another worker doesn't know how. Re-queue the task.
				// krancour: This behavior is something we can revisit in the future,
				// if and when we extract the async package into its own library.
				// Construct and execute a transaction that removes the task from this
				// worker's queue and re-queues it in the pending task queue.
				task.IncrementWorkerRejectionCount()
				newTaskJSON, err := task.ToJSON()
				if err != nil {
					select {
					case errCh <- fmt.Errorf(
						`error moving unprocessable task "%s" back to queue "%s": %s`,
						task.GetID(),
						pendingTaskQueueName,
						err,
					):
						continue
					case <-ctx.Done():
						return
					}
				}
				pipeline := e.redisClient.TxPipeline()
				pipeline.LPush(pendingTaskQueueName, newTaskJSON)
				pipeline.LRem(getActiveTaskQueueName(e.workerID), -1, taskJSON)
				_, err = pipeline.Exec()
				if err != nil {
					select {
					case errCh <- fmt.Errorf(
						`error moving unprocessable task "%s" back to queue "%s": %s`,
						task.GetID(),
						pendingTaskQueueName,
						err,
					):
					case <-ctx.Done():
						return
					}
				}
				continue
			}
			taskSuccess := false
			followUpTaskJSONs := [][]byte{}
			hadMarshalingError := false
			followUpTasks, err := jobFn(ctx, task)
			if err != nil {
				// If we get to here, we have a legitimate failure executing the task.
				// This isn't the worker's fault. Simply log this.
				// krancour: This behavior is something we can revisit in the future, if
				// and when we extract the async package into its own library.
				log.WithFields(log.Fields{
					"job":    task.GetJobName(),
					"taskID": task.GetID(),
					"error":  err,
				}).Error("error executing job; not submitting any follow-up tasks")
			} else {
				taskSuccess = true
				// We might have follow-up tasks that we need to enqueue. We should
				// do as much prep-work as we can for that BEFORE starting a
				// transaction. This way, if there's a failure, we can log it, and then
				// still, at least, try to execute a smaller transaction that JUST
				// removes the current task from the active task queue. This is because
				// we don't want this task getting STUCK in the active task queue where
				// a cleaner will eventually put it back on the pending task queue when
				// this worker dies.
				for _, followUpTask := range followUpTasks {
					// In reality, this is nearly guaranteed to never fail because there's
					// no legitimate possibility of a task not being serializable. So it's
					// possible that the following is unnecessarily defensive.
					followUpTaskJSON, err := followUpTask.ToJSON()
					if err != nil {
						hadMarshalingError = true
						log.WithFields(log.Fields{
							"job":            task.GetJobName(),
							"taskID":         task.GetID(),
							"followUpJob":    followUpTask.GetJobName(),
							"followUpTaskID": followUpTask.GetID(),
							"error":          err,
						}).Error(
							"error marshaling follow-up task; not submitting any follow-up " +
								"tasks",
						)
						// Don't break; continue. We want to log all the failures; not just
						// the first.
						continue
					}
					followUpTaskJSONs = append(followUpTaskJSONs, followUpTaskJSON)
				}
			}
			// Regardless of success or failure, we're done with this task. Remove it
			// from the active task queue.
			pipeline := e.redisClient.TxPipeline()
			pipeline.LRem(getActiveTaskQueueName(e.workerID), -1, taskJSON)
			// If the task was successful and we had no trouble marshaling the
			// follow-up tasks, we can add them to the appropriate queues
			if taskSuccess && !hadMarshalingError {
				for i, followUpTask := range followUpTasks {
					if followUpTask.GetExecuteTime() != nil {
						pipeline.LPush(deferredTaskQueueName, followUpTaskJSONs[i])
					} else {
						pipeline.LPush(pendingTaskQueueName, followUpTaskJSONs[i])
					}
				}
			}
			_, err = pipeline.Exec()
			if err != nil {
				// At this point, we're only possibly dealing with a Redis failure.
				// Unlike some of the conditions above, there's nothing we can do.
				// This is fatal.
				select {
				case errCh <- fmt.Errorf(
					`error removing task "%s" from queue "%s" and submitting follow-up `+
						`tasks: %s`,
					e.workerID,
					getActiveTaskQueueName(e.workerID),
					err,
				):
					continue
				case <-ctx.Done():
					return
				}
			}
		case <-ctx.Done():
			return
		}
	}
}

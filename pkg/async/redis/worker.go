package redis

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/Azure/open-service-broker-azure/pkg/async"
	log "github.com/Sirupsen/logrus"
	"github.com/go-redis/redis"
	uuid "github.com/satori/go.uuid"
)

// receiveTasksFn defines functions used to receive tasks from one queue and
// dispatch them to another
type receiveTasksFn func(
	ctx context.Context,
	sourceQueueName string,
	destinationQueueName string,
	retCh chan []byte,
	errCh chan error,
)

// executeTasksFn defines functions used to execute pending tasks
type executeTasksFn func(
	ctx context.Context,
	inputCh chan []byte,
	pendingTaskQueueName string,
	errCh chan error,
)

// watchDeferredTaskFn defines functions used to watch a deferred task
type watchDeferredTaskFn func(
	ctx context.Context,
	taskJSON []byte,
	pendingTaskQueueName string,
	errCh chan error,
)

// Worker is an interface to be implemented by components that receive and
// asynchronously complete provisioning and deprovisioning tasks
type Worker interface {
	// GetID returns the worker's ID
	GetID() string
	// RegisterJob registers a new Job with the worker
	RegisterJob(name string, fn async.JobFn) error
	// Run causes the worker to complete tasks. It blocks until a fatal error is
	// encountered or the context passed to it has been canceled. Run always
	// returns a non-nil error.
	Run(context.Context) error
}

// worker is a Redis-based implementation of the Worker interface
type worker struct {
	id          string
	redisClient *redis.Client
	// This allows tests to inject an alternative implementation
	heart        Heart
	jobsFns      map[string]async.JobFn
	jobsFnsMutex sync.RWMutex
	// This allows tests to inject an alternative implementation of this function
	receivePendingTasks receiveTasksFn
	// This allows tests to inject an alternative implementation of this function
	receiveDeferredTasks receiveTasksFn
	// This allows tests to inject an alternative implementation of this function
	executeTasks executeTasksFn
	// This allows tests to inject an alternative implementation of this function
	watchDeferredTask watchDeferredTaskFn
}

// newWorker returns a new Redis-based implementation of the Worker interface
func newWorker(redisClient *redis.Client) Worker {
	workerID := uuid.NewV4().String()
	w := &worker{
		id:          workerID,
		redisClient: redisClient,
		heart:       newHeart(workerID, time.Second*30, redisClient),
		jobsFns:     make(map[string]async.JobFn),
	}
	w.receivePendingTasks = w.defaultReceiveTasks
	w.receiveDeferredTasks = w.defaultReceiveTasks
	w.executeTasks = w.defaultExecuteTasks
	w.watchDeferredTask = w.defaultWatchDeferredTask
	return w
}

// GetID returns the worker's ID
func (w *worker) GetID() string {
	return w.id
}

// RegisterJob registers a new Job with the worker
func (w *worker) RegisterJob(name string, fn async.JobFn) error {
	w.jobsFnsMutex.Lock()
	defer w.jobsFnsMutex.Unlock()
	if _, ok := w.jobsFns[name]; ok {
		return &errDuplicateJob{name: name}
	}
	w.jobsFns[name] = fn
	return nil
}

// Run causes the worker to complete tasks. It blocks until a fatal error is
// encountered or the context passed to it has been canceled. Run always returns
// a non-nil error.
func (w *worker) Run(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	errCh := make(chan error)
	// As soon as we add the worker to the workers set, it's eligible for the
	// cleaner to clean up after it, so it's important that we guarantee the
	// cleaner will see this worker as alive. We can't trust that the heartbeat
	// loop (which we'll shortly start in its own goroutine) will have sent the
	// first heartbeat BEFORE the worker is added to the workers set. To account
	// for this, we synchronously send the first heartbeat.
	if err := w.heart.Beat(); err != nil {
		return err
	}
	// Heartbeat loop
	go func() {
		select {
		case errCh <- &errHeartStopped{workerID: w.id, err: w.heart.Run(ctx)}:
		case <-ctx.Done():
		}
	}()
	// Announce this worker's existence
	err := w.redisClient.SAdd(workerSetName, w.id).Err()
	if err != nil {
		return fmt.Errorf(
			`error adding worker "%s" to worker set: %s`,
			w.id,
			err,
		)
	}
	// Assemble and execute a pipeline to receive and execute pending tasks...
	go func() {
		pendingReceiverRetCh := make(chan []byte)
		pendingReceiverErrCh := make(chan error)
		executorErrCh := make(chan error)
		go w.receivePendingTasks(
			ctx,
			pendingTaskQueueName,
			getActiveTaskQueueName(w.id),
			pendingReceiverRetCh,
			pendingReceiverErrCh,
		)
		// Fan out to 5 executors
		for range [5]struct{}{} {
			go w.executeTasks(
				ctx,
				pendingReceiverRetCh,
				pendingTaskQueueName,
				executorErrCh,
			)
		}
		select {
		case err := <-pendingReceiverErrCh:
			errCh <- &errReceiverStopped{
				workerID:  w.id,
				queueName: pendingTaskQueueName,
				err:       err,
			}
		case err := <-executorErrCh:
			errCh <- &errTaskExecutorStopped{workerID: w.id, err: err}
		case <-ctx.Done():
		}
	}()
	// Assemble and execute a pipeline to receive and watch deferred tasks...
	go func() {
		deferredReceiverRetCh := make(chan []byte)
		deferredReceiverErrCh := make(chan error)
		watcherErrCh := make(chan error)
		go w.receiveDeferredTasks(
			ctx,
			deferredTaskQueueName,
			getWatchedTaskQueueName(w.id),
			deferredReceiverRetCh,
			deferredReceiverErrCh,
		)
		// Fan out to as many watchers as we need
		go func() {
			for {
				select {
				case taskJSON := <-deferredReceiverRetCh:
					w.watchDeferredTask(
						ctx,
						taskJSON,
						pendingTaskQueueName,
						watcherErrCh,
					)
				case <-ctx.Done():
					return
				}
			}
		}()
		select {
		case err := <-deferredReceiverErrCh:
			errCh <- &errReceiverStopped{
				workerID:  w.id,
				queueName: deferredTaskQueueName,
				err:       err,
			}
		case err := <-watcherErrCh:
			errCh <- &errDeferredTaskWatcherStopped{workerID: w.id, err: err}
		case <-ctx.Done():
		}
	}()
	// Now wait...
	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
		log.Debug("context canceled; async worker shutting down")
		return ctx.Err()
	}
}

// defaultReceive receives tasks from a source queue and dispatches them to a
// to both a destination queue and a return channel.
func (w *worker) defaultReceiveTasks(
	ctx context.Context,
	sourceQueueName string,
	destinationQueueName string,
	retCh chan []byte,
	errCh chan error,
) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	for {
		taskJSON, err := w.redisClient.BRPopLPush(
			sourceQueueName,
			destinationQueueName,
			time.Second*5,
		).Bytes()
		if err == redis.Nil {
			select {
			case <-ctx.Done():
				return
			default:
				continue
			}
		}
		if err != nil {
			select {
			case errCh <- fmt.Errorf(
				`error receiving task from queue "%s": %s`,
				sourceQueueName,
				err,
			):
				continue
			case <-ctx.Done():
				return
			}
		}
		select {
		case retCh <- taskJSON:
		case <-ctx.Done():
			return
		}
	}
}

func (w *worker) defaultExecuteTasks(
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
			task, err := w.getTaskFromJSON(taskJSON, getActiveTaskQueueName(w.id))
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
			w.jobsFnsMutex.RLock()
			defer w.jobsFnsMutex.RUnlock()
			jobFn, ok := w.jobsFns[task.GetJobName()]
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
				pipeline := w.redisClient.TxPipeline()
				pipeline.LPush(pendingTaskQueueName, newTaskJSON)
				pipeline.LRem(getActiveTaskQueueName(w.id), -1, taskJSON)
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
			if _, err := jobFn(ctx, task); err != nil {
				// If we get to here, we have a legitimate failure executing the task.
				// This isn't the worker's fault. Simply log this.
				// krancour: This behavior is something we can revisit in the future, if
				// and when we extract the async package into its own library.
				log.WithFields(log.Fields{
					"job":    task.GetJobName(),
					"taskID": task.GetID(),
					"error":  err,
				}).Error("error executing job")
			}
			// Regardless of success or failure, we're done with this task. Remove it
			// from the active work queue.
			err = w.redisClient.LRem(getActiveTaskQueueName(w.id), -1, taskJSON).Err()
			if err != nil {
				select {
				case errCh <- fmt.Errorf(
					`error removing task "%s" from queue "%s": %s`,
					w.id,
					getActiveTaskQueueName(w.id),
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

func (w *worker) getTaskFromJSON(
	taskJSON []byte,
	queueName string,
) (async.Task, error) {
	task, err := async.NewTaskFromJSON(taskJSON)
	if err != nil {
		// If the JSON is invalid, remove the message from the queue, log this and
		// move on. No other worker is going to be able to process this-- there's
		// nothing we can do and there's no sense treating this as a fatal
		// condition.
		err := w.redisClient.LRem(queueName, -1, taskJSON).Err()
		if err != nil {
			return nil, fmt.Errorf(
				`error removing malformed task from queue "%s"; task: %s: %s`,
				queueName,
				taskJSON,
				err,
			)
		}
		log.WithFields(log.Fields{
			"queue":    queueName,
			"taskJSON": taskJSON,
			"error":    err,
		}).Error("error decoding malformed task from queue")
		return nil, nil
	}
	return task, nil
}

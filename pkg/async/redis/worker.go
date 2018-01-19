package redis

import (
	"context"
	"fmt"
	"sync"

	"github.com/Azure/open-service-broker-azure/pkg/async"
	log "github.com/Sirupsen/logrus"
	"github.com/go-redis/redis"
	uuid "github.com/satori/go.uuid"
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
	id           string
	jobsFns      map[string]async.JobFn
	jobsFnsMutex sync.RWMutex
	redisClient  *redis.Client
	// This allows tests to inject an alternative implementation of this function
	runHeart runHeartFn
	// This allows tests to inject an alternative implementation of this function
	heartbeat heartbeatFn
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
		jobsFns:     make(map[string]async.JobFn),
	}
	w.runHeart = w.defaultRunHeart
	w.heartbeat = w.defaultHeartbeat
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
	if err := w.heartbeat(); err != nil {
		return err
	}
	// Heartbeat loop
	go func() {
		select {
		case errCh <- &errHeartStopped{workerID: w.id, err: w.runHeart(ctx)}:
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
				deferredTaskQueueName,
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

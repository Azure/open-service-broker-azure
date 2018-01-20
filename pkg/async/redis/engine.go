package redis

import (
	"context"
	"fmt"

	"github.com/Azure/open-service-broker-azure/pkg/async"
	log "github.com/Sirupsen/logrus"
	"github.com/go-redis/redis"
	uuid "github.com/satori/go.uuid"
)

// engine is a Redis-based implementation of the Engine interface.
type engine struct {
	redisClient *redis.Client
	workerID    string
	// This allows tests to inject an alternative implementation of this function
	clean cleanFn
	// This allows tests to inject an alternative implementation of this function
	cleanActiveTaskQueue cleanWorkerQueueFn
	// This allows tests to inject an alternative implementation of this function
	cleanWatchedTaskQueue cleanWorkerQueueFn
	// This allows tests to inject an alternative implementation of this function
	runHeart runHeartFn
	// This allows tests to inject an alternative implementation of this function
	heartbeat heartbeatFn
	// This allows tests to inject an alternative implementation of Worker
	worker Worker
}

// NewEngine returns a new Redis-based implementation of the aync.Engine
// interface
func NewEngine(redisClient *redis.Client) async.Engine {
	workerID := uuid.NewV4().String()
	e := &engine{
		workerID:    workerID,
		redisClient: redisClient,
		worker:      newWorker(redisClient, workerID),
	}
	e.clean = e.defaultClean
	e.cleanActiveTaskQueue = e.defaultCleanWorkerQueue
	e.cleanWatchedTaskQueue = e.defaultCleanWorkerQueue
	e.runHeart = e.defaultRunHeart
	e.heartbeat = e.defaultHeartbeat
	return e
}

// RegisterJob registers a new async.JobFn with the async engine
func (e *engine) RegisterJob(name string, fn async.JobFn) error {
	return e.worker.RegisterJob(name, fn)
}

// SubmitTask submits an idempotent task to the async engine for reliable,
// asynchronous completion
func (e *engine) SubmitTask(task async.Task) error {
	taskJSON, err := task.ToJSON()
	if err != nil {
		return fmt.Errorf("error encoding task %#v: %s", task, err)
	}

	var queueName string
	if task.GetExecuteTime() != nil {
		queueName = deferredTaskQueueName
	} else {
		queueName = pendingTaskQueueName
	}

	err = e.redisClient.LPush(queueName, taskJSON).Err()
	if err != nil {
		return fmt.Errorf("error encoding task %#v: %s", task, err)
	}
	return nil
}

// Run causes the async engine to carry out all of its functions. It blocks
// until a fatal error is encountered or the context passed to it has been
// canceled. Run always returns a non-nil error.
func (e *engine) Run(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	errCh := make(chan error)
	// Start the cleaner
	go func() {
		select {
		case errCh <- &errCleanerStopped{
			err: e.clean(
				ctx,
				workerSetName,
				pendingTaskQueueName,
				deferredTaskQueueName,
			),
		}:
		case <-ctx.Done():
		}
	}()
	// As soon as we add the worker to the workers set, it's eligible for the
	// cleaner to clean up after it, so it's important that we guarantee the
	// cleaner will see this worker as alive. We can't trust that the heartbeat
	// loop (which we'll shortly start in its own goroutine) will have sent the
	// first heartbeat BEFORE the worker is added to the workers set. To account
	// for this, we synchronously send the first heartbeat.
	if err := e.heartbeat(); err != nil {
		return err
	}
	// Heartbeat loop
	go func() {
		select {
		case errCh <- &errHeartStopped{workerID: e.workerID, err: e.runHeart(ctx)}:
		case <-ctx.Done():
		}
	}()
	// Announce this worker's existence
	err := e.redisClient.SAdd(workerSetName, e.workerID).Err()
	if err != nil {
		return fmt.Errorf(
			`error adding worker "%s" to worker set: %s`,
			e.workerID,
			err,
		)
	}
	// Start the worker
	go func() {
		select {
		case errCh <- &errWorkerStopped{
			workerID: e.worker.GetID(),
			err:      e.worker.Run(ctx),
		}:
		case <-ctx.Done():
		}
	}()
	select {
	case <-ctx.Done():
		log.Debug("context canceled; async engine shutting down")
		return ctx.Err()
	case err := <-errCh:
		return err
	}
}

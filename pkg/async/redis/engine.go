package redis

import (
	"context"
	"crypto/tls"
	"fmt"
	"sync"

	"github.com/Azure/open-service-broker-azure/pkg/async"
	log "github.com/Sirupsen/logrus"
	"github.com/go-redis/redis"
	uuid "github.com/satori/go.uuid"
)

// engine is a Redis-based implementation of the Engine interface.
type engine struct {
	workerID     string
	jobsFns      map[string]async.JobFn
	jobsFnsMutex sync.RWMutex
	redisClient  *redis.Client
	config       Config
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
	// This allows tests to inject an alternative implementation of this function
	receivePendingTasks receiveTasksFn
	// This allows tests to inject an alternative implementation of this function
	receiveDeferredTasks receiveTasksFn
	// This allows tests to inject an alternative implementation of this function
	executeTasks executeTasksFn
	// This allows tests to inject an alternative implementation of this function
	watchDeferredTasks watchDeferredTasksFn
}

// NewEngine returns a new Redis-based implementation of the aync.Engine
// interface
func NewEngine(config Config) async.Engine {
	workerID := uuid.NewV4().String()
	redisOpts := &redis.Options{
		Addr:       fmt.Sprintf("%s:%d", config.RedisHost, config.RedisPort),
		Password:   config.RedisPassword,
		DB:         config.RedisDB,
		MaxRetries: 5,
	}
	if config.RedisEnableTLS {
		redisOpts.TLSConfig = &tls.Config{
			ServerName: config.RedisHost,
		}
	}
	e := &engine{
		workerID:    workerID,
		jobsFns:     make(map[string]async.JobFn),
		redisClient: redis.NewClient(redisOpts),
		config:      config,
	}
	e.clean = e.defaultClean
	e.cleanActiveTaskQueue = e.defaultCleanWorkerQueue
	e.cleanWatchedTaskQueue = e.defaultCleanWorkerQueue
	e.runHeart = e.defaultRunHeart
	e.heartbeat = e.defaultHeartbeat
	e.receivePendingTasks = e.defaultReceiveTasks
	e.receiveDeferredTasks = e.defaultReceiveTasks
	e.executeTasks = e.defaultExecuteTasks
	e.watchDeferredTasks = e.defaultWatchDeferredTasks
	return e
}

// RegisterJob registers a new async.JobFn with the async engine
func (e *engine) RegisterJob(name string, fn async.JobFn) error {
	e.jobsFnsMutex.Lock()
	defer e.jobsFnsMutex.Unlock()
	if _, ok := e.jobsFns[name]; ok {
		return &errDuplicateJob{name: name}
	}
	e.jobsFns[name] = fn
	return nil
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
				cleaningInterval,
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
	if err := e.heartbeat(cleaningInterval * 2); err != nil {
		return err
	}
	// Heartbeat loop
	go func() {
		select {
		case errCh <- &errHeartStopped{
			workerID: e.workerID,
			err:      e.runHeart(ctx, cleaningInterval),
		}:
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
	// Assemble and execute a pipeline to receive and execute pending tasks...
	go func() {
		pendingReceiverRetCh := make(chan []byte)
		pendingReceiverErrCh := make(chan error)
		executorErrCh := make(chan error)
		go e.receivePendingTasks(
			ctx,
			pendingTaskQueueName,
			getActiveTaskQueueName(e.workerID),
			pendingReceiverRetCh,
			pendingReceiverErrCh,
		)
		// Fan out to desired number of workers
		for i := 0; i < e.config.PendingTaskWorkerCount; i++ {
			go e.executeTasks(
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
				workerID:  e.workerID,
				queueName: pendingTaskQueueName,
				err:       err,
			}
		case err := <-executorErrCh:
			errCh <- &errTaskExecutorStopped{workerID: e.workerID, err: err}
		case <-ctx.Done():
		}
	}()
	// Assemble and execute a pipeline to receive and watch deferred tasks...
	go func() {
		deferredReceiverRetCh := make(chan []byte)
		deferredReceiverErrCh := make(chan error)
		watcherErrCh := make(chan error)
		go e.receiveDeferredTasks(
			ctx,
			deferredTaskQueueName,
			getWatchedTaskQueueName(e.workerID),
			deferredReceiverRetCh,
			deferredReceiverErrCh,
		)
		// Fan out to desired number of watchers
		for i := 0; i < e.config.DeferedTaskWatcherCount; i++ {
			go e.watchDeferredTasks(
				ctx,
				deferredReceiverRetCh,
				pendingTaskQueueName,
				watcherErrCh,
			)
		}
		select {
		case err := <-deferredReceiverErrCh:
			errCh <- &errReceiverStopped{
				workerID:  e.workerID,
				queueName: deferredTaskQueueName,
				err:       err,
			}
		case err := <-watcherErrCh:
			errCh <- &errDeferredTaskWatcherStopped{workerID: e.workerID, err: err}
		case <-ctx.Done():
		}
	}()
	// Now wait...
	select {
	case <-ctx.Done():
		log.Debug("context canceled; async engine shutting down")
		return ctx.Err()
	case err := <-errCh:
		return err
	}
}

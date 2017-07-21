package async

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/Azure/azure-service-broker/pkg/async/model"
	"github.com/go-redis/redis"
	uuid "github.com/satori/go.uuid"
)

type receiveAndWorkFunction func(ctx context.Context, queueName string) error
type workFunction func(ctx context.Context, task model.Task) error

// Worker is an interface to be implemented by components that receive and
// asynchronously complete provisioning and deprovisioning tasks
type Worker interface {
	// GetID returns the worker's ID
	GetID() string
	// RegisterJob registers a new Job with the worker
	RegisterJob(name string, fn model.JobFunction) error
	// Work causes the worker to begin completing tasks
	Work(context.Context) error
}

// worker is a Redis-based implementation of the Worker interface
type worker struct {
	id          string
	redisClient *redis.Client
	// This allows tests to inject an alternative implementation
	heart        Heart
	jobsFns      map[string]model.JobFunction
	jobsFnsMutex sync.RWMutex
	// TODO: Split this into two functions for better testability
	// This allows tests to inject an alternative implementation of this function
	receiveAndWork receiveAndWorkFunction
	// This allows tests to inject an alternative implementation of this function
	work workFunction
}

// newWorker returns a new Reids-based implementation of the Worker interface
func newWorker(redisClient *redis.Client) Worker {
	workerID := uuid.NewV4().String()
	w := &worker{
		id:          workerID,
		redisClient: redisClient,
		heart:       newHeart(workerID, time.Second*30, redisClient),
		jobsFns:     make(map[string]model.JobFunction),
	}
	w.receiveAndWork = w.defaultReceiveAndWork
	w.work = w.defaultWork
	return w
}

// GetID returns the worker's ID
func (w *worker) GetID() string {
	return w.id
}

// RegisterJob registers a new Job with the worker
func (w *worker) RegisterJob(name string, fn model.JobFunction) error {
	w.jobsFnsMutex.Lock()
	defer w.jobsFnsMutex.Unlock()
	if _, ok := w.jobsFns[name]; ok {
		return &errDuplicateJob{name: name}
	}
	w.jobsFns[name] = fn
	return nil
}

// Work causes the worker to begin completing tasks
func (w *worker) Work(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	errChan := make(chan error)
	// As soon as we add the worker to the workers set, it's eligible for the
	// cleaner to clean up after it, so it's important that we guarantee the
	// cleaner will see this worker as alive. We can't trust that the heartbeat
	// loop (which we'll shortly start in its own goroutine) will have sent the
	// first heartbeat BEFORE the worker is added to the workers set. To account
	// for this, we synchronously send the first heartbeat.
	err := w.heart.Beat()
	if err != nil {
		return err
	}
	// Heartbeat loop
	go func() {
		err := w.heart.Start(ctx)
		hse := &errHeartStopped{workerID: w.id, err: err}
		select {
		case errChan <- hse:
		case <-ctx.Done():
		}
	}()
	// Announce this worker's existence
	intCmd := w.redisClient.SAdd("workers", w.id)
	if intCmd.Err() != nil {
		return fmt.Errorf(
			`error adding worker "%s" to worker set: %s`,
			w.id,
			intCmd.Err(),
		)
	}
	// Receive and do work
	for range [5]struct{}{} {
		go func() {
			err := w.receiveAndWork(ctx, mainWorkQueueName)
			rawse := &errReceiveAndWorkStopped{workerID: w.id, err: err}
			select {
			case errChan <- rawse:
			case <-ctx.Done():
			}
		}()
	}
	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-errChan:
		return err
	}
}

// defaultReceiveAndWork synchronously receives and completes work. By combining
// these two operations, a worker never receives more work than it currently
// has the capacity to process.
func (w *worker) defaultReceiveAndWork(
	ctx context.Context,
	queueName string,
) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	for {
		strCmd := w.redisClient.BRPopLPush(
			queueName,
			getWorkerQueueName(w.id),
			time.Second*5,
		)
		if strCmd.Err() != redis.Nil {
			if strCmd.Err() != nil {
				return fmt.Errorf("error receiving task: %s", strCmd.Err())
			}
			taskJSON, err := strCmd.Result()
			if err != nil {
				return fmt.Errorf("error receiving task: %s", err)
			}
			task, err := model.NewTaskFromJSONString(taskJSON)
			if err != nil {
				// If the JSON is invalid, remove the message from this worker's queue,
				// log this and move on. No other worker is going to be able to process
				// this-- there's notthing we can do and there's no sense letting this
				// whole process die over this.
				log.Printf("error decoding task %s: %s", taskJSON, err)
				intCmd := w.redisClient.LRem(
					getWorkerQueueName(w.id),
					0,
					taskJSON,
				)
				if intCmd.Err() != nil {
					return fmt.Errorf(
						"error removing malformed task from the worker's work queue; task: %s: %s",
						taskJSON,
						err,
					)
				}
				continue
			}
			if err := w.work(ctx, task); err != nil {
				if _, ok := err.(*errJobNotFound); ok {
					// The error is that this worker doesn't know how to process this
					// task. That doesn't mean another worker doesn't know how. Re-queue
					// the task.
					// NB: This behavior is something we can revisit in the future if and
					// when we extract the async package into its own library.
					// Construct and execute a transaction that removes the task from this
					// worker's queue and re-queues it in the main work queue.
					task.IncrementWorkerRejectionCount()
					newTaskJSON, err := task.ToJSONString()
					if err != nil {
						return fmt.Errorf(
							"error moving unprocessable task back to main work queue; task: %#v: %s",
							task,
							err,
						)
					}
					pipeline := w.redisClient.TxPipeline()
					pipeline.LPush(queueName, newTaskJSON)
					pipeline.LRem(
						getWorkerQueueName(w.id),
						0,
						taskJSON,
					)
					_, err = pipeline.Exec()
					if err != nil {
						return fmt.Errorf(
							"error moving unprocessable task back to main work queue; task: %#v: %s",
							task,
							err,
						)
					}
					continue
				}
				// If we get to here, we have a legitimate failure executing the task.
				// This isn't the worker's fault. Simply log this.
				// NB: This behavior is something we can revisit in the future if and
				// when we extract the async package into its own library.
				log.Printf(
					`error running job "%s" for task: %#v: %s`,
					task.GetJobName(),
					task,
					err,
				)
			}
			intCmd := w.redisClient.LRem(
				getWorkerQueueName(w.id),
				0,
				taskJSON,
			)
			if intCmd.Err() != nil {
				return fmt.Errorf(
					`error removing task %s from worker "%s" work queue: %s`,
					taskJSON,
					w.id,
					err,
				)
			}
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
	}
}

func (w *worker) defaultWork(ctx context.Context, task model.Task) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	w.jobsFnsMutex.RLock()
	defer w.jobsFnsMutex.RUnlock()
	jobFn, ok := w.jobsFns[task.GetJobName()]
	if !ok {
		return &errJobNotFound{name: task.GetJobName()}
	}
	return jobFn(ctx, task.GetArgs())
}

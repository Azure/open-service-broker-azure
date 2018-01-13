package async

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/Azure/open-service-broker-azure/pkg/async/model"
	log "github.com/Sirupsen/logrus"
	"github.com/go-redis/redis"
	uuid "github.com/satori/go.uuid"
)

type receiveAndWorkFunction func(ctx context.Context, queueName string) error
type handleDelayedTasksFunction func(
	ctx context.Context,
	delayedQueueName string,
	activeQueueName string,
) error
type workFunction func(ctx context.Context, task model.Task) error
type handleDelayedTaskFunction func(
	ctx context.Context,
	taskJSON []byte,
	activeQueueName string,
	errCh chan error,
)

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
	// This allows tests to inject an alternative implementation of this function
	receiveAndWork receiveAndWorkFunction
	// This allows tests to inject an alternative implementation of this function
	handleDelayedTasks handleDelayedTasksFunction
	// This allows tests to inject an alternative implementation of this function
	work workFunction
	// This allows tests to inject an alternative implementation of this function
	handleDelayedTask handleDelayedTaskFunction
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
	w.handleDelayedTasks = w.defaultHandleDelayedTasks
	w.work = w.defaultWork
	w.handleDelayedTask = w.defaultHandleDelayedTask
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
	if err := w.heart.Beat(); err != nil {
		return err
	}
	// Heartbeat loop
	go func() {
		select {
		case errChan <- &errHeartStopped{workerID: w.id, err: w.heart.Start(ctx)}:
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
			select {
			case errChan <- &errReceiveAndWorkStopped{
				workerID: w.id,
				err:      w.receiveAndWork(ctx, mainActiveWorkQueueName),
			}:
			case <-ctx.Done():
			}
		}()
	}
	// Handle delayed tasks
	go func() {
		select {
		case errChan <- &errWatchDelayedTasksStopped{
			workerID: w.id,
			err: w.handleDelayedTasks(
				ctx,
				mainDelayedWorkQueueName,
				mainActiveWorkQueueName,
			),
		}:
		case <-ctx.Done():
		}
	}()
	select {
	case <-ctx.Done():
		log.Debug("context canceled; async worker shutting down")
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
			getWorkerActiveQueueName(w.id),
			time.Second*5,
		)
		// If we actually got something...
		if strCmd.Err() != redis.Nil {
			if strCmd.Err() != nil {
				return fmt.Errorf("error receiving active task: %s", strCmd.Err())
			}
			taskJSON, err := strCmd.Bytes()
			if err != nil {
				return fmt.Errorf("error receiving active task: %s", err)
			}
			task, err := model.NewTaskFromJSON(taskJSON)
			if err != nil {
				// If the JSON is invalid, remove the message from this worker's queue,
				// log this and move on. No other worker is going to be able to process
				// this-- there's nothing we can do and there's no sense letting this
				// whole process die over this.
				log.WithFields(log.Fields{
					"taskJSON": taskJSON,
					"error":    err,
				}).Error("error decoding active task")
				intCmd := w.redisClient.LRem(
					getWorkerActiveQueueName(w.id),
					0,
					taskJSON,
				)
				if intCmd.Err() != nil {
					return fmt.Errorf(
						"error removing malformed task from the worker's active work "+
							"queue; task: %s: %s",
						taskJSON,
						intCmd.Err(),
					)
				}
				continue
			}
			if err := w.work(ctx, task); err != nil {
				if _, ok := err.(*errJobNotFound); ok {
					// The error is that this worker doesn't know how to process this
					// task. That doesn't mean another worker doesn't know how. Re-queue
					// the task.
					// krancour: This behavior is something we can revisit in the future
					// if and when we extract the async package into its own library.
					// Construct and execute a transaction that removes the task from this
					// worker's queue and re-queues it in the main work queue.
					task.IncrementWorkerRejectionCount()
					newTaskJSON, err := task.ToJSON()
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
						getWorkerActiveQueueName(w.id),
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
				// krancour: This behavior is something we can revisit in the future if
				// and when we extract the async package into its own library.
				log.WithFields(log.Fields{
					"job":    task.GetJobName(),
					"taskID": task.GetID(),
					"error":  err,
				}).Error("error executing job")
			}
			intCmd := w.redisClient.LRem(
				getWorkerActiveQueueName(w.id),
				0,
				taskJSON,
			)
			if intCmd.Err() != nil {
				return fmt.Errorf(
					`error removing task %s from worker "%s" work queue: %s`,
					taskJSON,
					w.id,
					intCmd.Err(),
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

func (w *worker) defaultHandleDelayedTasks(
	ctx context.Context,
	delayedQueueName string,
	activeQueueName string,
) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	handleDelayedTaskErrCh := make(chan error)
	for {
		strCmd := w.redisClient.BRPopLPush(
			delayedQueueName,
			getWorkerDelayedQueueName(w.id),
			time.Second*5,
		)
		// If we actually got something...
		if strCmd.Err() != redis.Nil {
			if strCmd.Err() != nil {
				return fmt.Errorf("error receiving delayed task: %s", strCmd.Err())
			}
			taskJSON, err := strCmd.Bytes()
			if err != nil {
				return fmt.Errorf("error receiving delayed task: %s", err)
			}
			// Launch a new goroutine to deal with it.
			// krancour: Goroutines are very lightweight and it is practical to have
			// even hundreds of thousands in a single process. We use a finite number
			// of goroutines to do actual work because we cannot accurately forecast
			// the resource requirements of the tasks that are being executed, but
			// if ALL we need to do in this case is wait for a certain wall-clock time
			// to roll around, a goroutine for handling that is very cheap.
			// Practically speaking, we can afford to spawn as many of those as we
			// like.
			go w.handleDelayedTask(
				ctx,
				taskJSON,
				activeQueueName,
				handleDelayedTaskErrCh,
			)
		}
		select {
		case err := <-handleDelayedTaskErrCh:
			return err
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

func (w *worker) defaultHandleDelayedTask(
	ctx context.Context,
	taskJSON []byte,
	mainActiveWorkQueueName string,
	errCh chan error,
) {
	task, err := model.NewTaskFromJSON(taskJSON)
	if err != nil {
		// If the JSON is invalid, remove the message from this worker's queue,
		// log this and move on. No other worker is going to be able to process
		// this-- there's nothing we can do and there's no sense letting this
		// whole process die over this.
		log.WithFields(log.Fields{
			"taskJSON": taskJSON,
			"error":    err,
		}).Error("error decoding delayed task")
		intCmd := w.redisClient.LRem(
			getWorkerDelayedQueueName(w.id),
			0,
			taskJSON,
		)
		if intCmd.Err() != nil {
			errCh <- fmt.Errorf(
				"error removing malformed task from the worker's delayed work "+
					"queue; task: %s: %s",
				taskJSON,
				intCmd.Err(),
			)
		}
		return
	}
	executeTime := task.GetExecuteTime()
	if executeTime == nil {
		intCmd := w.redisClient.LRem(
			getWorkerDelayedQueueName(w.id),
			0,
			taskJSON,
		)
		if intCmd.Err() != nil {
			errCh <- fmt.Errorf(
				"error removing task with no executeTime from the worker's delayed "+
					"work queue; task: %s: %s",
				taskJSON,
				intCmd.Err(),
			)
			return
		}
		log.WithFields(log.Fields{
			"task": task.GetID(),
		}).Error("delayed task had no executeTime and was removed from the " +
			"delayed work queue")
		return
	}
	// Note if the duration passed to the timer is 0 or negative, it should go
	// off immediately
	timer := time.NewTimer(time.Until(*executeTime))
	defer timer.Stop()
	select {
	case <-timer.C:
		pipeline := w.redisClient.TxPipeline()
		pipeline.LPush(mainActiveWorkQueueName, taskJSON)
		pipeline.LRem(
			getWorkerDelayedQueueName(w.id),
			0,
			taskJSON,
		)
		_, err = pipeline.Exec()
		if err != nil {
			errCh <- fmt.Errorf(
				`error moving delayed task "%s" back to main work queue: %s`,
				task.GetID(),
				err,
			)
		}
	case <-ctx.Done():
		errCh <- ctx.Err()
	}
}

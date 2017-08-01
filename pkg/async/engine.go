package async

import (
	"context"
	"fmt"

	"github.com/Azure/azure-service-broker/pkg/async/model"
	log "github.com/Sirupsen/logrus"
	"github.com/go-redis/redis"
)

// Engine is an interface for a broker-specifc framework for submitting and
// asynchronously completing provisioning and deprovisioning tasks.
type Engine interface {
	// RegisterJob registers a new Job with the async engine
	RegisterJob(name string, fn model.JobFunction) error
	// SubmitTask submits an idempotent task to the async engine for reliable,
	// asynchronous completion
	SubmitTask(model.Task) error
	// Start causes the async engine to begin executing queued tasks
	Start(context.Context) error
}

// engine is a Redis-based implementation of the Engine interface.
type engine struct {
	redisClient *redis.Client
	// This allows tests to inject an alternative implementation of Worker
	worker Worker
	// This allows tests to inject an alternative implementation of Cleaner
	cleaner Cleaner
}

// NewEngine returns a new Redis-based implementation of the Engine
// interface
func NewEngine(redisClient *redis.Client) Engine {
	return &engine{
		redisClient: redisClient,
		cleaner:     newCleaner(redisClient),
		worker:      newWorker(redisClient),
	}
}

// RegisterJob registers a new Job with the async engine
func (e *engine) RegisterJob(name string, fn model.JobFunction) error {
	return e.worker.RegisterJob(name, fn)
}

// SubmitTask submits an idempotent task to the async engine for reliable,
// asynchronous completion
func (e *engine) SubmitTask(task model.Task) error {
	taskJSON, err := task.ToJSON()
	if err != nil {
		return fmt.Errorf("error encoding task %#v: %s", task, err)
	}
	intCmd := e.redisClient.LPush(mainWorkQueueName, taskJSON)
	if intCmd.Err() != nil {
		return fmt.Errorf("error encoding task %#v: %s", task, err)
	}
	return nil
}

func (e *engine) Start(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	errChan := make(chan error)
	// Start the cleaner
	go func() {
		select {
		case errChan <- &errCleanerStopped{err: e.cleaner.Clean(ctx)}:
		case <-ctx.Done():
		}
	}()
	// Start the worker
	go func() {
		select {
		case errChan <- &errWorkerStopped{
			workerID: e.worker.GetID(),
			err:      e.worker.Work(ctx),
		}:
		case <-ctx.Done():
		}
	}()
	select {
	case <-ctx.Done():
		log.Debug("context canceled; async engine shutting down")
		return ctx.Err()
	case err := <-errChan:
		return err
	}
}

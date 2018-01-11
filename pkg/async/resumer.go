package async

import (
	"context"
	"fmt"
	"sync"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/go-redis/redis"
)

type resumeFunction func(mainWorkQueueName string) error

// Resumer is an interface to be implemented by components that
// resume delayed tasks on a periodic cadence
type Resumer interface {
	Resume(context.Context) error
	Watch(queue string)
}

type resumer struct {
	redisClient       *redis.Client
	resume            resumeFunction
	watchedQueues     map[string]struct{}
	watchedQueueMutex sync.Mutex
}

func newResumer(redisClient *redis.Client) Resumer {
	r := &resumer{
		redisClient:   redisClient,
		watchedQueues: make(map[string]struct{}),
	}
	r.resume = r.defaultResumeFunc
	return r
}

// Resume moves delayed tasks from watched queues to
// the main worker queue
func (r *resumer) Resume(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	ticker := time.NewTicker(time.Minute * 5)
	for {
		if err := r.resume(mainWorkQueueName); err != nil {
			return &errResuming{err: err}
		}
		select {
		case <-ticker.C:
		case <-ctx.Done():
			log.Debug("context canceled; async task resumer shutting down")
			return ctx.Err()
		}
	}
}

// Watch adds a queue to the set of watched queues
func (r *resumer) Watch(queueName string) {
	r.watchedQueueMutex.Lock()
	defer r.watchedQueueMutex.Unlock()
	r.watchedQueues[queueName] = struct{}{}
}

func (r *resumer) defaultResumeFunc(mainWorkQueueName string) error {
	r.watchedQueueMutex.Lock()
	defer r.watchedQueueMutex.Unlock()
	for queueName := range r.watchedQueues {
		strCmd := r.redisClient.RPopLPush(queueName, mainWorkQueueName)
		for strCmd.Err() != redis.Nil {
			if strCmd.Err() != nil {
				log.WithFields(log.Fields{
					"error": strCmd.Err(),
				}).Fatal("error starting delayed task")
				return fmt.Errorf("error starting delayed task %s", strCmd.Err())
			}
			strCmd = r.redisClient.RPopLPush(queueName, mainWorkQueueName)
		}
		//Should be empty now, delete it from the map. It will
		//get readded if another is task is delayed
		delete(r.watchedQueues, queueName)
	}
	return nil
}

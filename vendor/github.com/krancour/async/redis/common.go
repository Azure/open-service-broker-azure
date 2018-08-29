package redis

import (
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/krancour/async"
)

const (
	workerSetName         = "workers"
	aliveIndicator        = "alive"
	pendingTaskQueueName  = "pendingTasks"
	deferredTaskQueueName = "deferredTasks"
)

func getActiveTaskQueueName(workerID string) string {
	return fmt.Sprintf("active-tasks:%s", workerID)
}

func getWatchedTaskQueueName(workerID string) string {
	return fmt.Sprintf("watched-tasks:%s", workerID)
}

func (e *engine) getTaskFromJSON(
	taskJSON []byte,
	queueName string,
) (async.Task, error) {
	task, err := async.NewTaskFromJSON(taskJSON)
	if err != nil {
		// If the JSON is invalid, remove the message from the queue, log this and
		// move on. No other worker is going to be able to process this-- there's
		// nothing we can do and there's no sense treating this as a fatal
		// condition.
		err := e.redisClient.LRem(queueName, -1, taskJSON).Err()
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

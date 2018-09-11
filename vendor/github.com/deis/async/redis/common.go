package redis

import (
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/deis/async"
)

const (
	aliveIndicator = "alive"
)

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

func (e *engine) prefixRedisKey(key string) string {
	if e.config.RedisPrefix != "" {
		return fmt.Sprintf("%s:%s", e.config.RedisPrefix, key)
	}
	return key
}

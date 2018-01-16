package redis

import (
	"errors"

	goredis "github.com/go-redis/redis"
	uuid "github.com/satori/go.uuid"
)

var (
	redisClient = goredis.NewClient(&goredis.Options{
		Addr:     "redis:6379",
		PoolSize: 20,
	})
	errSome = errors.New("an error")
)

func getDisposableQueueName() string {
	return uuid.NewV4().String()
}

func getDisposableWorkerID() string {
	return uuid.NewV4().String()
}

func getDisposableWorkerSetName() string {
	return uuid.NewV4().String()
}

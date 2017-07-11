package async

import (
	"errors"

	"github.com/go-redis/redis"
)

var (
	redisClient = redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		PoolSize: 20,
	})
	errSome = errors.New("an error")
)

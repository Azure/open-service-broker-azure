package main

import (
	"fmt"
	"log"

	"context"

	"github.com/go-redis/redis"
)

func main() {
	redisConfig, err := getRedisConfig()
	if err != nil {
		log.Fatal(err)
	}
	redisClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", redisConfig.Host, redisConfig.Port),
		Password: redisConfig.Password,
		DB:       redisConfig.DB,
	})
	broker, err := newBroker(redisClient, modules)
	if err != nil {
		log.Fatal(err)
	}
	if err := broker.start(context.Background()); err != nil {
		log.Fatal(err)
	}
}

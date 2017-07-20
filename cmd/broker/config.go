package main

import (
	"github.com/kelseyhightower/envconfig"
)

// redisConfig represents details for connecting to the Redis instance that
// the broker itself relies on for storing state and orchestrating asynchronous
// processes
type redisConfig struct {
	Host     string `envconfig:"REDIS_HOST" required:"true"`
	Port     int    `envconfig:"REDIS_PORT" default:"6379"`
	Password string `envconfig:"REDIS_PASSWORD" default:""`
	DB       int    `envconfig:"REDIS_DB" default:"0"`
}

func getRedisConfig() (redisConfig, error) {
	redisConfig := redisConfig{}
	err := envconfig.Process("", &redisConfig)
	return redisConfig, err
}

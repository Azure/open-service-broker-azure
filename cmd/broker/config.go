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

// cryptoConfig represents details (e.g. key) for encrypting and decrypting any
// (potentially) sensitive information
type cryptoConfig struct {
	AES256Key string `envconfig:"AES256_KEY" required:"true"`
}

func getRedisConfig() (redisConfig, error) {
	redisConfig := redisConfig{}
	err := envconfig.Process("", &redisConfig)
	return redisConfig, err
}

func getCryptoConfig() (cryptoConfig, error) {
	cryptoConfig := cryptoConfig{}
	err := envconfig.Process("", &cryptoConfig)
	return cryptoConfig, err
}

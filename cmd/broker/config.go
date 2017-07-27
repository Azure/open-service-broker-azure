package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/kelseyhightower/envconfig"
)

// logConfig represents configuration options for the broker's leveled logging
type logConfig struct {
	LevelStr string `envconfig:"LOG_LEVEL" default:"INFO"`
	Level    log.Level
}

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

type basicAuthConfig struct {
	Username string `envconfig:"BASIC_AUTH_USERNAME" required:"true"`
	Password string `envconfig:"BASIC_AUTH_PASSWORD" required:"true"`
}

func getLogConfig() (logConfig, error) {
	logConfig := logConfig{}
	err := envconfig.Process("", &logConfig)
	if err != nil {
		return logConfig, err
	}
	logConfig.Level, err = log.ParseLevel(logConfig.LevelStr)
	return logConfig, err
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

func getBasicAuthConfig() (basicAuthConfig, error) {
	basicAuthConfig := basicAuthConfig{}
	err := envconfig.Process("", &basicAuthConfig)
	return basicAuthConfig, err
}

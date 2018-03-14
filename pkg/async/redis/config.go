package redis

import "github.com/kelseyhightower/envconfig"

const envconfigPrefix = "ASYNC"

// Config encapsulates all configuration options for the Redis-based
// implementation of the async.Engine interface
type Config struct {
	RedisHost               string `envconfig:"REDIS_HOST" required:"true"`
	RedisPort               int    `envconfig:"REDIS_PORT"`
	RedisPassword           string `envconfig:"REDIS_PASSWORD"`
	RedisDB                 int    `envconfig:"REDIS_DB"`
	RedisEnableTLS          bool   `envconfig:"REDIS_ENABLE_TLS"`
	PendingTaskWorkerCount  int    `envconfig:"PENDING_TASK_WORKER_COUNT"`
	DeferedTaskWatcherCount int    `envconfig:"DEFERED_TASK_WATCHER_COUNT"`
}

// NewConfigWithDefaults returns a Config object with default values already
// applied. Callers are then free to set custom values for the remaining fields
// and/or override default values.
func NewConfigWithDefaults() Config {
	return Config{
		RedisPort:               6379,
		PendingTaskWorkerCount:  5,
		DeferedTaskWatcherCount: 100,
	}
}

// GetConfigFromEnvironment returns configuration derived from environment
// variables
func GetConfigFromEnvironment() (Config, error) {
	c := NewConfigWithDefaults()
	err := envconfig.Process(envconfigPrefix, &c)
	return c, err
}

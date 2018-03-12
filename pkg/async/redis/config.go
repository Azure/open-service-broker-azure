package redis

import "github.com/kelseyhightower/envconfig"

// Config encapsulates all configuration options for the Redis-based
// implementation of the async.Engine interface
type Config struct {
	RedisHost      string `envconfig:"ASYNC_REDIS_HOST" required:"true"`
	RedisPort      int    `envconfig:"ASYNC_REDIS_PORT"`
	RedisPassword  string `envconfig:"ASYNC_REDIS_PASSWORD"`
	RedisDB        int    `envconfig:"ASYNC_REDIS_DB"`
	RedisEnableTLS bool   `envconfig:"ASYNC_REDIS_ENABLE_TLS"`
}

// NewConfigWithDefaults returns a Config object with default values already
// applied. Callers are then free to set custom values for the remaining fields
// and/or override default values.
func NewConfigWithDefaults() Config {
	return Config{RedisPort: 6379}
}

// GetConfigFromEnvironment returns configuration derived from environment
// variables
func GetConfigFromEnvironment() (Config, error) {
	c := NewConfigWithDefaults()
	err := envconfig.Process("", &c)
	return c, err
}

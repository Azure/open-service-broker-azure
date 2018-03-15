package storage

import "github.com/kelseyhightower/envconfig"

// Config represents configuration options for the Redis-based implementation
// of the Store interface
type Config struct {
	RedisHost      string `envconfig:"STORAGE_REDIS_HOST" required:"true"`
	RedisPort      int    `envconfig:"STORAGE_REDIS_PORT"`
	RedisPassword  string `envconfig:"STORAGE_REDIS_PASSWORD"`
	RedisDB        int    `envconfig:"STORAGE_REDIS_DB"`
	RedisEnableTLS bool   `envconfig:"STORAGE_REDIS_ENABLE_TLS"`
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

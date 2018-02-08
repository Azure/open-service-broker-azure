package config

import "github.com/kelseyhightower/envconfig"

// RedisConfig represents details for connecting to the Redis instance that
// the broker itself relies on for storing state and orchestrating asynchronous
// processes
type RedisConfig interface {
	GetHost() string
	GetPort() int
	GetPassword() string
	GetStorageDB() int
	GetAsyncDB() int
	IsTLSEnabled() bool
}

type redisConfig struct {
	Host      string `envconfig:"REDIS_HOST" required:"true"`
	Port      int    `envconfig:"REDIS_PORT" default:"6379"`
	Password  string `envconfig:"REDIS_PASSWORD" default:""`
	StorageDB int    `envconfig:"REDIS_STORAGE_DB" default:"0"`
	AsyncDB   int    `envconfig:"REDIS_ASYNC_DB" default:"1"`
	EnableTLS bool   `envconfig:"REDIS_ENABLE_TLS" default:"false"`
}

// GetRedisConfig returns Redis configuration
func GetRedisConfig() (RedisConfig, error) {
	rc := redisConfig{}
	err := envconfig.Process("", &rc)
	return rc, err
}

func (r redisConfig) GetHost() string {
	return r.Host
}

func (r redisConfig) GetPort() int {
	return r.Port
}

func (r redisConfig) GetPassword() string {
	return r.Password
}

func (r redisConfig) GetStorageDB() int {
	return r.StorageDB
}

func (r redisConfig) GetAsyncDB() int {
	return r.AsyncDB
}

func (r redisConfig) IsTLSEnabled() bool {
	return r.EnableTLS
}

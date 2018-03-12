package storage

import "github.com/kelseyhightower/envconfig"

// RedisConfig represents details for connecting to the Redis instance that
// the broker relies on for storage
type RedisConfig interface {
	GetHost() string
	GetPort() int
	GetPassword() string
	GetDB() int
	IsTLSEnabled() bool
}

type redisConfig struct {
	Host      string `envconfig:"STORAGE_REDIS_HOST" required:"true"`
	Port      int    `envconfig:"STORAGE_REDIS_PORT" default:"6379"`
	Password  string `envconfig:"STORAGE_REDIS_PASSWORD" default:""`
	DB        int    `envconfig:"STORAGE_REDIS_DB" default:"0"`
	EnableTLS bool   `envconfig:"STORAGE_REDIS_ENABLE_TLS" default:"false"`
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

func (r redisConfig) GetDB() int {
	return r.DB
}

func (r redisConfig) IsTLSEnabled() bool {
	return r.EnableTLS
}

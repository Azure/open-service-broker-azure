package redis

import "github.com/kelseyhightower/envconfig"

// Config represents details for connecting to the Redis instance that
// the broker relies on for orchestrating asynchronous processes
type Config interface {
	GetHost() string
	GetPort() int
	GetPassword() string
	GetDB() int
	IsTLSEnabled() bool
}

type config struct {
	Host      string `envconfig:"ASYNC_REDIS_HOST" required:"true"`
	Port      int    `envconfig:"ASYNC_REDIS_PORT" default:"6379"`
	Password  string `envconfig:"ASYNC_REDIS_PASSWORD" default:""`
	DB        int    `envconfig:"ASYNC_REDIS_DB" default:"0"`
	EnableTLS bool   `envconfig:"ASYNC_REDIS_ENABLE_TLS" default:"false"`
}

// GetConfig returns Redis configuration
func GetConfig() (Config, error) {
	c := config{}
	err := envconfig.Process("", &c)
	return c, err
}

func (c config) GetHost() string {
	return c.Host
}

func (c config) GetPort() int {
	return c.Port
}

func (c config) GetPassword() string {
	return c.Password
}

func (c config) GetDB() int {
	return c.DB
}

func (c config) IsTLSEnabled() bool {
	return c.EnableTLS
}

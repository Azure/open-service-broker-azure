package redis

import (
	"strings"

	"github.com/kelseyhightower/envconfig"
)

const envconfigPrefix = "STORAGE"

// Config represents configuration options for the Redis-based implementation
// of the Store interface
type Config struct {
	RedisHost        string `envconfig:"REDIS_HOST" required:"true"`
	RedisPort        int    `envconfig:"REDIS_PORT"`
	RedisPassword    string `envconfig:"REDIS_PASSWORD"`
	RedisDB          int    `envconfig:"REDIS_DB"`
	RedisEnableTLS   bool   `envconfig:"REDIS_ENABLE_TLS"`
	EncryptionScheme string `envconfig:"ENCRYPTION_SCHEME" default:"AES256"`
	AES256Key        string `envconfig:"AES256_KEY"`
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
	err := envconfig.Process(envconfigPrefix, &c)
	if err != nil {
		return c, err
	}
	c.EncryptionScheme = strings.ToUpper(c.EncryptionScheme)
	return c, nil
}

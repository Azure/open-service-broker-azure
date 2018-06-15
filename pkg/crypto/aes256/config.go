package aes256

import (
	"github.com/kelseyhightower/envconfig"
)

const envconfigPrefix = "CRYPTO"

// Config represents configuration options for the AES256-based implementation
// of the Crypto interface
type Config struct {
	Key string `envconfig:"AES256_KEY"`
}

// NewConfigWithDefaults returns a Config object with default values already
// applied. Callers are then free to set custom values for the remaining fields
// and/or override default values.
func NewConfigWithDefaults() Config {
	return Config{}
}

// GetConfigFromEnvironment returns configuration derived from environment
// variables
func GetConfigFromEnvironment() (Config, error) {
	c := NewConfigWithDefaults()
	err := envconfig.Process(envconfigPrefix, &c)
	return c, err
}

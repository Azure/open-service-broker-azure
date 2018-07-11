package api

import (
	"github.com/kelseyhightower/envconfig"
)

const envconfigPrefix = "API_SERVER"

// Config represents configuration options for the API server
type Config struct {
	Port        int    `envconfig:"PORT"`
	TLSCertPath string `envconfig:"TLS_CERT_PATH"`
	TLSKeyPath  string `envconfig:"TLS_KEY_PATH"`
}

// NewConfigWithDefaults returns a Config object with default values already
// applied. Callers are then free to set custom values for the remaining fields
// and/or override default values.
func NewConfigWithDefaults() Config {
	return Config{Port: 8080}
}

// GetConfigFromEnvironment returns configuration derived from environment
// variables
func GetConfigFromEnvironment() (Config, error) {
	c := NewConfigWithDefaults()
	err := envconfig.Process(envconfigPrefix, &c)
	return c, err
}

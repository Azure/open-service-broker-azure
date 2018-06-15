package crypto

import (
	"strings"

	"github.com/kelseyhightower/envconfig"
)

const envconfigPrefix = "CRYPTO"

// Config represents configuration options for the global codec
type Config struct {
	EncryptionScheme string `envconfig:"ENCRYPTION_SCHEME" default:"AES256"`
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
	if err != nil {
		return c, err
	}
	c.EncryptionScheme = strings.ToUpper(c.EncryptionScheme)
	return c, nil
}

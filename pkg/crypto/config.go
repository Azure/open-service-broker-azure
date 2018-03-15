package crypto

import (
	"errors"
	"fmt"
	"strings"

	"github.com/kelseyhightower/envconfig"
)

// Config represents details (e.g. key) for encrypting and decrypting any
// (potentially) sensitive information
type Config interface {
	GetEncryptionScheme() string
	GetAES256Key() string
}

type config struct {
	EncryptionScheme string `envconfig:"ENCRYPTION_SCHEME" default:"AES256"`
	AES256Key        string `envconfig:"AES256_KEY"`
}

// GetConfig returns crypto configuration
func GetConfig() (Config, error) {
	cc := config{}
	err := envconfig.Process("", &cc)
	cc.EncryptionScheme = strings.ToUpper(cc.EncryptionScheme)
	switch cc.EncryptionScheme {
	case AES256:
		if cc.AES256Key == "" {
			return cc, errors.New("AES256_KEY was not specified")
		}
		if len(cc.AES256Key) != 32 {
			return cc, errors.New("AES256_KEY is an invalid length")
		}
	case NOOP:
	default:
		return cc, fmt.Errorf(
			`unrecognized ENCRYPTION_SCHEME "%s"`,
			cc.EncryptionScheme,
		)
	}
	return cc, err
}

func (c config) GetEncryptionScheme() string {
	return c.EncryptionScheme
}

func (c config) GetAES256Key() string {
	return c.AES256Key
}

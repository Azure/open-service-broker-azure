package config

import (
	"errors"
	"fmt"
	"strings"

	"github.com/Azure/open-service-broker-azure/pkg/crypto"
	"github.com/kelseyhightower/envconfig"
)

// CryptoConfig represents details (e.g. key) for encrypting and decrypting any
// (potentially) sensitive information
type CryptoConfig interface {
	GetEncryptionScheme() string
	GetAES256Key() string
}

type cryptoConfig struct {
	EncryptionScheme string `envconfig:"ENCRYPTION_SCHEME" default:"AES256"`
	AES256Key        string `envconfig:"AES256_KEY"`
}

// GetCryptoConfig returns crypto configuration
func GetCryptoConfig() (CryptoConfig, error) {
	cc := cryptoConfig{}
	err := envconfig.Process("", &cc)
	cc.EncryptionScheme = strings.ToUpper(cc.EncryptionScheme)
	switch cc.EncryptionScheme {
	case crypto.AES256:
		if cc.AES256Key == "" {
			return cc, errors.New("AES256_KEY was not specified")
		}
		if len(cc.AES256Key) != 32 {
			return cc, errors.New("AES256_KEY is an invalid length")
		}
	case crypto.NOOP:
	default:
		return cc, fmt.Errorf(
			`unrecognized ENCRYPTION_SCHEME "%s"`,
			cc.EncryptionScheme,
		)
	}
	return cc, err
}

func (c cryptoConfig) GetEncryptionScheme() string {
	return c.EncryptionScheme
}

func (c cryptoConfig) GetAES256Key() string {
	return c.AES256Key
}

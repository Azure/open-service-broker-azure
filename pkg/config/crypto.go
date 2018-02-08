package config

import "github.com/kelseyhightower/envconfig"

// CryptoConfig represents details (e.g. key) for encrypting and decrypting any
// (potentially) sensitive information
type CryptoConfig interface {
	GetAES256Key() string
}

type cryptoConfig struct {
	AES256Key string `envconfig:"AES256_KEY" required:"true"`
}

// GetCryptoConfig returns crypto configuration
func GetCryptoConfig() (CryptoConfig, error) {
	cc := cryptoConfig{}
	err := envconfig.Process("", &cc)
	return cc, err
}

func (c cryptoConfig) GetAES256Key() string {
	return c.AES256Key
}

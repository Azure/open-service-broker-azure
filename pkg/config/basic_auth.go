package config

import "github.com/kelseyhightower/envconfig"

// BasicAuthConfig represents details such as username and password that will
// be used to secure the broker using basic auth
type BasicAuthConfig interface {
	GetUsername() string
	GetPassword() string
}

type basicAuthConfig struct {
	Username string `envconfig:"BASIC_AUTH_USERNAME" required:"true"`
	Password string `envconfig:"BASIC_AUTH_PASSWORD" required:"true"`
}

// GetBasicAuthConfig returns basic auth configuration
func GetBasicAuthConfig() (BasicAuthConfig, error) {
	bac := basicAuthConfig{}
	err := envconfig.Process("", &bac)
	return bac, err
}

func (b basicAuthConfig) GetUsername() string {
	return b.Username
}

func (b basicAuthConfig) GetPassword() string {
	return b.Password
}

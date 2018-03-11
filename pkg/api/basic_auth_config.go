package api

import (
	"github.com/kelseyhightower/envconfig"
)

// BasicAuthConfig represents details such as username and password that will
// be used to secure the broker using basic auth
type BasicAuthConfig interface {
	GetUsername() string
	GetPassword() string
}

type basicAuthConfig struct {
	Username   string `envconfig:"BASIC_AUTH_USERNAME"`
	CFUsername string `envconfig:"SECURITY_USER_NAME"`
	Password   string `envconfig:"BASIC_AUTH_PASSWORD"`
	CFPassword string `envconfig:"SECURITY_USER_PASSWORD"`
}

// GetBasicAuthConfig returns basic auth configuration
func GetBasicAuthConfig() (BasicAuthConfig, error) {
	bac := basicAuthConfig{}
	err := envconfig.Process("", &bac)
	if bac.Username == "" && bac.CFUsername == "" {
		return bac, &errBasicAuthUsernameNotSpecified{}
	}
	if bac.Username == "" && bac.CFUsername != "" {
		bac.Username = bac.CFUsername
	}
	if bac.Password == "" && bac.CFPassword == "" {
		return bac, &errBasicAuthPasswordNotSpecified{}
	}
	if bac.Password == "" && bac.CFPassword != "" {
		bac.Password = bac.CFPassword
	}
	return bac, err
}

func (b basicAuthConfig) GetUsername() string {
	return b.Username
}

func (b basicAuthConfig) GetPassword() string {
	return b.Password
}

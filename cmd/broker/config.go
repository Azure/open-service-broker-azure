package main

import (
	"fmt"
	"strings"

	"github.com/Azure/open-service-broker-azure/pkg/service"
	log "github.com/Sirupsen/logrus"
	"github.com/kelseyhightower/envconfig"
)

// logConfig represents configuration options for the broker's leveled logging
type logConfig struct {
	LevelStr string `envconfig:"LOG_LEVEL" default:"INFO"`
	Level    log.Level
}

// redisConfig represents details for connecting to the Redis instance that
// the broker itself relies on for storing state and orchestrating asynchronous
// processes
type redisConfig struct {
	Host      string `envconfig:"REDIS_HOST" required:"true"`
	Port      int    `envconfig:"REDIS_PORT" default:"6379"`
	Password  string `envconfig:"REDIS_PASSWORD" default:""`
	DB        int    `envconfig:"REDIS_DB" default:"0"`
	EnableTLS bool   `envconfig:"REDIS_ENABLE_TLS" default:"false"`
}

// cryptoConfig represents details (e.g. key) for encrypting and decrypting any
// (potentially) sensitive information
type cryptoConfig struct {
	AES256Key string `envconfig:"AES256_KEY" required:"true"`
}

type basicAuthConfig struct {
	Username string `envconfig:"BASIC_AUTH_USERNAME" required:"true"`
	Password string `envconfig:"BASIC_AUTH_PASSWORD" required:"true"`
}

type modulesConfig struct {
	MinStabilityStr string `envconfig:"MIN_STABILITY" default:"EXPERIMENTAL"`
	MinStability    service.Stability
}

type azureConfig struct {
	DefaultLocation      string `envconfig:"AZURE_DEFAULT_LOCATION"`
	DefaultResourceGroup string `envconfig:"AZURE_DEFAULT_RESOURCE_GROUP"`
}

func getLogConfig() (logConfig, error) {
	lc := logConfig{}
	err := envconfig.Process("", &lc)
	if err != nil {
		return lc, err
	}
	lc.Level, err = log.ParseLevel(lc.LevelStr)
	return lc, err
}

func getRedisConfig() (redisConfig, error) {
	rc := redisConfig{}
	err := envconfig.Process("", &rc)
	return rc, err
}

func getCryptoConfig() (cryptoConfig, error) {
	cc := cryptoConfig{}
	err := envconfig.Process("", &cc)
	return cc, err
}

func getBasicAuthConfig() (basicAuthConfig, error) {
	bac := basicAuthConfig{}
	err := envconfig.Process("", &bac)
	return bac, err
}

func getModulesConfig() (modulesConfig, error) {
	mc := modulesConfig{}
	err := envconfig.Process("", &mc)
	if err != nil {
		return mc, err
	}
	minStabilityStr := strings.ToUpper(mc.MinStabilityStr)
	switch minStabilityStr {
	case "EXPERIMENTAL":
		mc.MinStability = service.StabilityExperimental
	case "PREVIEW":
		mc.MinStability = service.StabilityPreview
	case "STABLE":
		mc.MinStability = service.StabilityStable
	default:
		return mc, fmt.Errorf(
			`unrecognized stability level "%s"`,
			minStabilityStr,
		)
	}
	return mc, nil
}

func getAzureConfig() (azureConfig, error) {
	ac := azureConfig{}
	err := envconfig.Process("", &ac)
	return ac, err
}

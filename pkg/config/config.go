package config

import (
	"fmt"
	"strings"

	"github.com/Azure/open-service-broker-azure/pkg/service"
	log "github.com/Sirupsen/logrus"
	"github.com/kelseyhightower/envconfig"
)

// LogConfig represents configuration options for the broker's leveled logging
type LogConfig struct {
	LevelStr string `envconfig:"LOG_LEVEL" default:"INFO"`
	Level    log.Level
}

// RedisConfig represents details for connecting to the Redis instance that
// the broker itself relies on for storing state and orchestrating asynchronous
// processes
type RedisConfig struct {
	Host      string `envconfig:"REDIS_HOST" required:"true"`
	Port      int    `envconfig:"REDIS_PORT" default:"6379"`
	Password  string `envconfig:"REDIS_PASSWORD" default:""`
	StorageDB int    `envconfig:"REDIS_STORAGE_DB" default:"0"`
	AsyncDB   int    `envconfig:"REDIS_ASYNC_DB" default:"1"`
	EnableTLS bool   `envconfig:"REDIS_ENABLE_TLS" default:"false"`
}

// CryptoConfig represents details (e.g. key) for encrypting and decrypting any
// (potentially) sensitive information
type CryptoConfig struct {
	AES256Key string `envconfig:"AES256_KEY" required:"true"`
}

// BasicAuthConfig represents details such as username and password that will
// be used to secure the broker using basic auth
type BasicAuthConfig struct {
	Username string `envconfig:"BASIC_AUTH_USERNAME" required:"true"`
	Password string `envconfig:"BASIC_AUTH_PASSWORD" required:"true"`
}

// ModulesConfig represents details re: which modules should be included or
// excluded when the broker is started
type ModulesConfig struct {
	MinStabilityStr string `envconfig:"MIN_STABILITY" default:"EXPERIMENTAL"`
	MinStability    service.Stability
}

// AzureConfig represents details necessary for the broker to interact with
// an Azure subscription
type AzureConfig struct {
	Environment          string `envconfig:"AZURE_ENVIRONMENT" default:"AzurePublicCloud"` // nolint: lll
	SubscriptionID       string `envconfig:"AZURE_SUBSCRIPTION_ID" required:"true"`        // nolint: lll
	TenantID             string `envconfig:"AZURE_TENANT_ID" required:"true"`
	ClientID             string `envconfig:"AZURE_CLIENT_ID" required:"true"`
	ClientSecret         string `envconfig:"AZURE_CLIENT_SECRET" required:"true"`
	DefaultLocation      string `envconfig:"AZURE_DEFAULT_LOCATION"`
	DefaultResourceGroup string `envconfig:"AZURE_DEFAULT_RESOURCE_GROUP"`
}

// GetLogConfig returns log configuration
func GetLogConfig() (LogConfig, error) {
	lc := LogConfig{}
	err := envconfig.Process("", &lc)
	if err != nil {
		return lc, err
	}
	lc.Level, err = log.ParseLevel(lc.LevelStr)
	return lc, err
}

// GetRedisConfig returns Redis configuration
func GetRedisConfig() (RedisConfig, error) {
	rc := RedisConfig{}
	err := envconfig.Process("", &rc)
	return rc, err
}

// GetCryptoConfig returns crypto configuration
func GetCryptoConfig() (CryptoConfig, error) {
	cc := CryptoConfig{}
	err := envconfig.Process("", &cc)
	return cc, err
}

// GetBasicAuthConfig returns basic auth configuration
func GetBasicAuthConfig() (BasicAuthConfig, error) {
	bac := BasicAuthConfig{}
	err := envconfig.Process("", &bac)
	return bac, err
}

// GetModulesConfig returns modules configuration
func GetModulesConfig() (ModulesConfig, error) {
	mc := ModulesConfig{}
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

// GetAzureConfig returns Azure subscription configuration
func GetAzureConfig() (AzureConfig, error) {
	ac := AzureConfig{}
	err := envconfig.Process("", &ac)
	return ac, err
}

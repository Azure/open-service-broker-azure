package config

import (
	"fmt"
	"strings"

	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/open-service-broker-azure/pkg/service"
	"github.com/kelseyhightower/envconfig"
)

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
	EnvironmentStr       string `envconfig:"AZURE_ENVIRONMENT" default:"AzurePublicCloud"` // nolint: lll
	Environment          azure.Environment
	SubscriptionID       string `envconfig:"AZURE_SUBSCRIPTION_ID" required:"true"` // nolint: lll
	TenantID             string `envconfig:"AZURE_TENANT_ID" required:"true"`
	ClientID             string `envconfig:"AZURE_CLIENT_ID" required:"true"`
	ClientSecret         string `envconfig:"AZURE_CLIENT_SECRET" required:"true"`
	DefaultLocation      string `envconfig:"AZURE_DEFAULT_LOCATION"`
	DefaultResourceGroup string `envconfig:"AZURE_DEFAULT_RESOURCE_GROUP"`
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
	if err != nil {
		return ac, err
	}
	ac.Environment, err = azure.EnvironmentFromName(ac.EnvironmentStr)
	return ac, err
}

package config

import (
	"fmt"
	"strings"

	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/open-service-broker-azure/pkg/service"
	"github.com/kelseyhightower/envconfig"
)

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

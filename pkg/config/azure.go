package config

import (
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/kelseyhightower/envconfig"
)

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

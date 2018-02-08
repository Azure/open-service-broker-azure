package config

import (
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/kelseyhightower/envconfig"
)

// AzureConfig represents details necessary for the broker to interact with
// an Azure subscription
type AzureConfig interface {
	GetEnvironment() azure.Environment
	GetSubscriptionID() string
	GetTenantID() string
	GetClientID() string
	GetClientSecret() string
	GetDefaultLocation() string
	GetDefaultResourceGroup() string
}

type azureConfig struct {
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
	ac := azureConfig{}
	err := envconfig.Process("", &ac)
	if err != nil {
		return ac, err
	}
	ac.Environment, err = azure.EnvironmentFromName(ac.EnvironmentStr)
	return ac, err
}

func (a azureConfig) GetEnvironment() azure.Environment {
	return a.Environment
}

func (a azureConfig) GetSubscriptionID() string {
	return a.SubscriptionID
}

func (a azureConfig) GetTenantID() string {
	return a.TenantID
}

func (a azureConfig) GetClientID() string {
	return a.ClientID
}

func (a azureConfig) GetClientSecret() string {
	return a.ClientSecret
}

func (a azureConfig) GetDefaultLocation() string {
	return a.DefaultLocation
}

func (a azureConfig) GetDefaultResourceGroup() string {
	return a.DefaultResourceGroup
}

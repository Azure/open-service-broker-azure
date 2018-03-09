package azure

import (
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/kelseyhightower/envconfig"
)

// Config represents details necessary for the broker to interact with
// an Azure subscription
type Config interface {
	GetEnvironment() azure.Environment
	GetSubscriptionID() string
	GetTenantID() string
	GetClientID() string
	GetClientSecret() string
	GetDefaultLocation() string
	GetDefaultResourceGroup() string
}

type config struct {
	EnvironmentStr       string `envconfig:"AZURE_ENVIRONMENT" default:"AzurePublicCloud"` // nolint: lll
	Environment          azure.Environment
	SubscriptionID       string `envconfig:"AZURE_SUBSCRIPTION_ID" required:"true"` // nolint: lll
	TenantID             string `envconfig:"AZURE_TENANT_ID" required:"true"`
	ClientID             string `envconfig:"AZURE_CLIENT_ID" required:"true"`
	ClientSecret         string `envconfig:"AZURE_CLIENT_SECRET" required:"true"`
	DefaultLocation      string `envconfig:"AZURE_DEFAULT_LOCATION"`
	DefaultResourceGroup string `envconfig:"AZURE_DEFAULT_RESOURCE_GROUP"`
}

// GetConfig returns Azure-related configuration
func GetConfig() (Config, error) {
	c := config{}
	err := envconfig.Process("", &c)
	if err != nil {
		return c, err
	}
	c.Environment, err = azure.EnvironmentFromName(c.EnvironmentStr)
	return c, err
}

func (c config) GetEnvironment() azure.Environment {
	return c.Environment
}

func (c config) GetSubscriptionID() string {
	return c.SubscriptionID
}

func (c config) GetTenantID() string {
	return c.TenantID
}

func (c config) GetClientID() string {
	return c.ClientID
}

func (c config) GetClientSecret() string {
	return c.ClientSecret
}

func (c config) GetDefaultLocation() string {
	return c.DefaultLocation
}

func (c config) GetDefaultResourceGroup() string {
	return c.DefaultResourceGroup
}

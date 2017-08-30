package azure

import (
	"github.com/kelseyhightower/envconfig"
)

// Config represents all configuration details needed for connecting to Azure
// APIs
type Config struct {
	Environment    string `envconfig:"AZURE_ENVIRONMENT" default:"AzurePublicCloud"` // nolint: lll
	SubscriptionID string `envconfig:"AZURE_SUBSCRIPTION_ID" required:"true"`
	TenantID       string `envconfig:"AZURE_TENANT_ID" required:"true"`
	ClientID       string `envconfig:"AZURE_CLIENT_ID" required:"true"`
	ClientSecret   string `envconfig:"AZURE_CLIENT_SECRET" required:"true"`
}

// GetConfig parses configuration details needed for connecting to Azure APIs
// from environment variables and returns a Config object that encapsulates
// those details
func GetConfig() (Config, error) {
	ac := Config{}
	err := envconfig.Process("", &ac)
	return ac, err
}

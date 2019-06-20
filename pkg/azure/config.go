package azure

import (
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/kelseyhightower/envconfig"
)

const envconfigPrefix = "AZURE"

// Config represents details necessary for the broker to interact with
// an Azure subscription
type Config struct {
	Environment    azure.Environment
	SubscriptionID string `envconfig:"SUBSCRIPTION_ID" required:"true"`
	TenantID       string `envconfig:"TENANT_ID" required:"true"`
	ClientID       string `envconfig:"CLIENT_ID" required:"true"`
	ClientSecret   string `envconfig:"CLIENT_SECRET" required:"true"`
}

type tempConfig struct {
	Config
	EnvironmentStr string `envconfig:"ENVIRONMENT" default:"AzurePublicCloud"`
}

// NewConfigWithDefaults returns a Config object with default values already
// applied. Callers are then free to set custom values for the remaining fields
// and/or override default values.
func NewConfigWithDefaults() Config {
	return Config{}
}

// GetConfigFromEnvironment returns Azure-related configuration derived from
// environment variables
func GetConfigFromEnvironment() (Config, error) {
	c := tempConfig{
		Config: NewConfigWithDefaults(),
	}
	err := envconfig.Process(envconfigPrefix, &c)
	if err != nil {
		return c.Config, err
	}
	c.Environment, err = azure.EnvironmentFromName(c.EnvironmentStr)
	return c.Config, err
}

// GetEnvrionmentName returns the name of cloud envrionment.
// Expected return vaules are: ["AzurePublicCloud", "AzureChinaCloud"]
func GetEnvrionmentName() string {
	// We can directly ignore returned err here,
	// because this function is invoked at the start of
	// OSBA initiating. If there is an error to
	// invoke this function, OSBA will fail earlier.
	config, _ := GetConfigFromEnvironment()
	return config.Environment.Name
}

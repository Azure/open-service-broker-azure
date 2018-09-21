// +build !unit

package lifecycle

import (
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/open-service-broker-azure/pkg/azure"
)

func getAzureConfig() (
	*azure.Config,
	error,
) {
	azureConfig, err := azure.GetConfigFromEnvironment()
	if err != nil {
		return nil, err
	}
	return &azureConfig, nil
}

func getBearerTokenAuthorizer(azureConfig *azure.Config) (
	*autorest.BearerAuthorizer,
	error,
) {
	authorizer, err := azure.GetBearerTokenAuthorizer(
		azureConfig.Environment,
		azureConfig.TenantID,
		azureConfig.ClientID,
		azureConfig.ClientSecret,
	)
	if err != nil {
		return nil, err
	}

	return authorizer, nil
}

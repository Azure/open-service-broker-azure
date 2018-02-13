// +build !unit

package lifecycle

import (
	keyVaultSDK "github.com/Azure/azure-sdk-for-go/services/keyvault/mgmt/2016-10-01/keyvault" // nolint: lll
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/open-service-broker-azure/pkg/azure/arm"
	"github.com/Azure/open-service-broker-azure/pkg/config"
	"github.com/Azure/open-service-broker-azure/pkg/services/keyvault"
)

func getKeyvaultCases(
	azureEnvironment azure.Environment,
	subscriptionID string,
	authorizer autorest.Authorizer,
	armDeployer arm.Deployer,
) ([]serviceLifecycleTestCase, error) {
	azureConfig, err := config.GetAzureConfig()
	if err != nil {
		return nil, err
	}
	vaultsClient := keyVaultSDK.NewVaultsClientWithBaseURI(
		azureEnvironment.ResourceManagerEndpoint,
		azureConfig.GetSubscriptionID(),
	)
	vaultsClient.Authorizer = authorizer
	return []serviceLifecycleTestCase{
		{
			module: keyvault.New(
				azureConfig.GetTenantID(),
				armDeployer,
				vaultsClient,
			),
			serviceID: "d90c881e-c9bb-4e07-a87b-fcfe87e03276",
			planID:    "3577ee4a-75fc-44b3-b354-9d33d52ef486",
			location:  "southcentralus",
			provisioningParameters: &keyvault.ProvisioningParameters{
				ObjectID:     "6a74d229-e927-42c5-b6e8-8f5c095cfba8",
				ClientID:     "test",
				ClientSecret: "test",
			},
			bindingParameters: &keyvault.BindingParameters{},
		},
	}, nil
}

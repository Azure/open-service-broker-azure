// +build !unit

package lifecycle

import (
	storageSDK "github.com/Azure/azure-sdk-for-go/arm/storage"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/open-service-broker-azure/pkg/azure/arm"
	"github.com/Azure/open-service-broker-azure/pkg/services/storage"
)

func getStorageCases(
	azureEnvironment azure.Environment,
	subscriptionID string,
	authorizer autorest.Authorizer,
	armDeployer arm.Deployer,
) ([]serviceLifecycleTestCase, error) {
	accountsClient := storageSDK.NewAccountsClientWithBaseURI(
		azureEnvironment.ResourceManagerEndpoint,
		subscriptionID,
	)
	accountsClient.Authorizer = authorizer
	module := storage.New(armDeployer, accountsClient)
	return []serviceLifecycleTestCase{
		{ // General Purpose Storage Account
			module:                 module,
			description:            "general purpose storage account",
			serviceID:              "2e2fc314-37b6-4587-8127-8f9ee8b33fea",
			planID:                 "6ddf6b41-fb60-4b70-af99-8ecc4896b3cf",
			location:               "southcentralus",
			provisioningParameters: &storage.ProvisioningParameters{},
			bindingParameters:      &storage.BindingParameters{},
		},
		{ // Blob Storage Account
			module:                 module,
			description:            "blob storage account",
			serviceID:              "2e2fc314-37b6-4587-8127-8f9ee8b33fea",
			planID:                 "800a17e1-f20a-463d-a290-20516052f647",
			location:               "southcentralus",
			provisioningParameters: &storage.ProvisioningParameters{},
			bindingParameters:      &storage.BindingParameters{},
		},
		{ // Blob Storage Account + Blob Container
			module:                 module,
			description:            "blob storage account with a blob container",
			serviceID:              "2e2fc314-37b6-4587-8127-8f9ee8b33fea",
			planID:                 "189d3b8f-8307-4b3f-8c74-03d069237f70",
			location:               "southcentralus",
			provisioningParameters: &storage.ProvisioningParameters{},
			bindingParameters:      &storage.BindingParameters{},
		},
	}, nil
}

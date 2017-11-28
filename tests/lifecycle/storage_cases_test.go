// +build !unit

package lifecycle

import (
	"github.com/Azure/azure-service-broker/pkg/azure/arm"
	sa "github.com/Azure/azure-service-broker/pkg/azure/storage"
	"github.com/Azure/azure-service-broker/pkg/service"
	"github.com/Azure/azure-service-broker/pkg/services/storage"
)

func getStorageCases(
	armDeployer arm.Deployer,
	resourceGroup string,
) ([]moduleLifecycleTestCase, error) {
	storageManager, err := sa.NewManager()
	if err != nil {
		return nil, err
	}

	return []moduleLifecycleTestCase{
		{ // General Purpose Storage Account
			module:      storage.New(armDeployer, storageManager),
			description: "general purpose storage account",
			serviceID:   "2e2fc314-37b6-4587-8127-8f9ee8b33fea",
			planID:      "6ddf6b41-fb60-4b70-af99-8ecc4896b3cf",
			standardProvisioningContext: service.StandardProvisioningContext{
				Location: "southcentralus",
			},
			provisioningParameters: &storage.ProvisioningParameters{},
			bindingParameters:      &storage.BindingParameters{},
		},
		{ // Blob Storage Account
			module:      storage.New(armDeployer, storageManager),
			description: "blob storage account",
			serviceID:   "2e2fc314-37b6-4587-8127-8f9ee8b33fea",
			planID:      "800a17e1-f20a-463d-a290-20516052f647",
			standardProvisioningContext: service.StandardProvisioningContext{
				Location: "southcentralus",
			},
			provisioningParameters: &storage.ProvisioningParameters{},
			bindingParameters:      &storage.BindingParameters{},
		},
		{ // Blob Storage Account + Blob Container
			module:      storage.New(armDeployer, storageManager),
			description: "blob storage account with a blob container",
			serviceID:   "2e2fc314-37b6-4587-8127-8f9ee8b33fea",
			planID:      "189d3b8f-8307-4b3f-8c74-03d069237f70",
			standardProvisioningContext: service.StandardProvisioningContext{
				Location: "southcentralus",
			},
			provisioningParameters: &storage.ProvisioningParameters{},
			bindingParameters:      &storage.BindingParameters{},
		},
	}, nil
}

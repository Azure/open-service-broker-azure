// +build !unit

package lifecycle

import (
	"github.com/Azure/open-service-broker-azure/pkg/azure/arm"
	kv "github.com/Azure/open-service-broker-azure/pkg/azure/keyvault"
	"github.com/Azure/open-service-broker-azure/pkg/service"
	"github.com/Azure/open-service-broker-azure/pkg/services/keyvault"
)

func getKeyvaultCases(
	armDeployer arm.Deployer,
	resourceGroup string,
) ([]moduleLifecycleTestCase, error) {
	keyvaultManager, err := kv.NewManager()
	if err != nil {
		return nil, err
	}

	return []moduleLifecycleTestCase{
		{
			module:    keyvault.New(armDeployer, keyvaultManager),
			serviceID: "d90c881e-c9bb-4e07-a87b-fcfe87e03276",
			planID:    "3577ee4a-75fc-44b3-b354-9d33d52ef486",
			standardProvisioningContext: service.StandardProvisioningContext{
				Location: "southcentralus",
			},
			provisioningParameters: &keyvault.ProvisioningParameters{
				ObjectID:     "6a74d229-e927-42c5-b6e8-8f5c095cfba8",
				ClientID:     "test",
				ClientSecret: "test",
			},
			bindingParameters: &keyvault.BindingParameters{},
		},
	}, nil
}

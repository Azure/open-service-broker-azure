// +build !unit

package lifecycle

import (
	"github.com/Azure/azure-service-broker/pkg/azure/arm"
	eh "github.com/Azure/azure-service-broker/pkg/azure/eventhub"
	"github.com/Azure/azure-service-broker/pkg/services/eventhub"
)

func getEventhubCases(
	armDeployer arm.Deployer,
) ([]moduleLifecycleTestCase, error) {
	eventHubManager, err := eh.NewManager()
	if err != nil {
		return nil, err
	}

	return []moduleLifecycleTestCase{
		{
			module:    eventhub.New(armDeployer, eventHubManager),
			serviceID: "7bade660-32f1-4fd7-b9e6-d416d975170b",
			planID:    "80756db5-a20c-495d-ae70-62cf7d196a3c",
			provisioningParameters: &eventhub.ProvisioningParameters{
				Location:      "southcentralus",
				ResourceGroup: newTestResourceGroupName(),
			},
			bindingParameters: &eventhub.BindingParameters{},
		},
	}, nil
}

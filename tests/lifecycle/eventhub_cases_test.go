// +build !unit

package lifecycle

import (
	"github.com/Azure/open-service-broker-azure/pkg/azure/arm"
	eh "github.com/Azure/open-service-broker-azure/pkg/azure/eventhub"
	"github.com/Azure/open-service-broker-azure/pkg/service"
	"github.com/Azure/open-service-broker-azure/pkg/services/eventhubs"
)

func getEventhubCases(
	armDeployer arm.Deployer,
	resourceGroup string,
) ([]serviceLifecycleTestCase, error) {
	eventHubManager, err := eh.NewManager()
	if err != nil {
		return nil, err
	}

	return []serviceLifecycleTestCase{
		{
			module:    eventhubs.New(armDeployer, eventHubManager),
			serviceID: "7bade660-32f1-4fd7-b9e6-d416d975170b",
			planID:    "80756db5-a20c-495d-ae70-62cf7d196a3c",
			standardProvisioningContext: service.StandardProvisioningContext{
				Location: "southcentralus",
			},
			provisioningParameters: &eventhubs.ProvisioningParameters{},
			bindingParameters:      &eventhubs.BindingParameters{},
		},
	}, nil
}

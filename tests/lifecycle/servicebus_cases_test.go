// +build !unit

package lifecycle

import (
	"github.com/Azure/azure-service-broker/pkg/azure/arm"
	sb "github.com/Azure/azure-service-broker/pkg/azure/servicebus"
	"github.com/Azure/azure-service-broker/pkg/services/servicebus"
)

func getServicebusCases(
	armDeployer arm.Deployer,
) ([]moduleLifecycleTestCase, error) {
	serviceBusManager, err := sb.NewManager()
	if err != nil {
		return nil, err
	}

	return []moduleLifecycleTestCase{
		{
			module:    servicebus.New(armDeployer, serviceBusManager),
			serviceID: "6dc44338-2f13-4bc5-9247-5b1b3c5462d3",
			planID:    "13c6da8f-128c-48c0-a3a9-659d1b6d3920",
			provisioningParameters: &servicebus.ProvisioningParameters{
				Location:      "southcentralus",
				ResourceGroup: newTestResourceGroupName(),
			},
			bindingParameters: &servicebus.BindingParameters{},
		},
	}, nil
}

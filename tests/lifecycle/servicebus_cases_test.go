// +build !unit

package lifecycle

import (
	"github.com/Azure/azure-service-broker/pkg/azure/arm"
	sb "github.com/Azure/azure-service-broker/pkg/azure/servicebus"
	"github.com/Azure/azure-service-broker/pkg/service"
	"github.com/Azure/azure-service-broker/pkg/services/servicebus"
)

func getServicebusCases(
	armDeployer arm.Deployer,
	resourceGroup string,
) ([]moduleLifecycleTestCase, error) {
	serviceBusManager, err := sb.NewManager()
	if err != nil {
		return nil, err
	}

	return []moduleLifecycleTestCase{
		{
			module:    servicebus.New(armDeployer, serviceBusManager),
			serviceID: "6dc44338-2f13-4bc5-9247-5b1b3c5462d3",
			planID:    "d06817b1-87ea-4320-8942-14b1d060206a",
			standardProvisioningContext: service.StandardProvisioningContext{
				Location: "southcentralus",
			},
			provisioningParameters: &servicebus.ProvisioningParameters{},
			bindingParameters:      &servicebus.BindingParameters{},
		},
	}, nil
}

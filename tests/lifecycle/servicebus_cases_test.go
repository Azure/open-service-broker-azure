// +build !unit

package lifecycle

import (
	"github.com/Azure/open-service-broker-azure/pkg/azure/arm"
	sb "github.com/Azure/open-service-broker-azure/pkg/azure/servicebus"
	"github.com/Azure/open-service-broker-azure/pkg/service"
	"github.com/Azure/open-service-broker-azure/pkg/services/servicebus"
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
			planID:    "13c6da8f-128c-48c0-a3a9-659d1b6d3920",
			standardProvisioningContext: service.StandardProvisioningContext{
				Location: "southcentralus",
			},
			provisioningParameters: &servicebus.ProvisioningParameters{},
			bindingParameters:      &servicebus.BindingParameters{},
		},
	}, nil
}

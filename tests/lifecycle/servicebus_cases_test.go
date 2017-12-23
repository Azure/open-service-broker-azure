// +build !unit

package lifecycle

import (
	"github.com/Azure/open-service-broker-azure/pkg/azure/arm"
	sb "github.com/Azure/open-service-broker-azure/pkg/azure/servicebus"
	"github.com/Azure/open-service-broker-azure/pkg/services/servicebus"
)

func getServicebusCases(
	armDeployer arm.Deployer,
	resourceGroup string,
) ([]serviceLifecycleTestCase, error) {
	serviceBusManager, err := sb.NewManager()
	if err != nil {
		return nil, err
	}

	return []serviceLifecycleTestCase{
		{
			module:                 servicebus.New(armDeployer, serviceBusManager),
			serviceID:              "6dc44338-2f13-4bc5-9247-5b1b3c5462d3",
			planID:                 "d06817b1-87ea-4320-8942-14b1d060206a",
			location:               "southcentralus",
			provisioningParameters: &servicebus.ProvisioningParameters{},
			bindingParameters:      &servicebus.BindingParameters{},
		},
	}, nil
}

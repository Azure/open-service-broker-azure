// +build !unit

package lifecycle

import (
	"github.com/Azure/open-service-broker-azure/pkg/azure/arm"
	cr "github.com/Azure/open-service-broker-azure/pkg/azure/rediscache"
	"github.com/Azure/open-service-broker-azure/pkg/service"
	"github.com/Azure/open-service-broker-azure/pkg/services/acr"
)

func getAcrCases(
	armDeployer arm.Deployer,
	resourceGroup string,
) ([]moduleLifecycleTestCase, error) {
	acrManager, err := cr.NewManager()
	if err != nil {
		return nil, err
	}

	return []moduleLifecycleTestCase{
		{
			module:    acr.New(armDeployer, acrManager),
			serviceID: "0b9401d6-4c04-4b10-a4da-0fd6cd1c7b4a",
			planID:    "e74c0b8e-23b1-4495-91ae-41682d3a0b7c",
			standardProvisioningContext: service.StandardProvisioningContext{
				Location: "southcentralus",
			},
			provisioningParameters: &acr.ProvisioningParameters{},
			bindingParameters:      &acr.BindingParameters{},
		},
	}, nil
}

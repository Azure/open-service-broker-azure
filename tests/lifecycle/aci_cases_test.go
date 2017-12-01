// +build !unit

package lifecycle

import (
	ac "github.com/Azure/open-service-broker-azure/pkg/azure/aci"
	"github.com/Azure/open-service-broker-azure/pkg/azure/arm"
	"github.com/Azure/open-service-broker-azure/pkg/service"
	"github.com/Azure/open-service-broker-azure/pkg/services/aci"
)

func getACICases(
	armDeployer arm.Deployer,
	resourceGroup string,
) ([]moduleLifecycleTestCase, error) {
	aciManager, err := ac.NewManager()
	if err != nil {
		return nil, err
	}

	return []moduleLifecycleTestCase{
		{
			module:    aci.New(armDeployer, aciManager),
			serviceID: "451d5d19-4575-4d4a-9474-116f705ecc95",
			planID:    "d48798e2-21db-405b-abc7-aa6f0ff08f6c",
			standardProvisioningContext: service.StandardProvisioningContext{
				Location: "eastus",
			},
			provisioningParameters: &aci.ProvisioningParameters{
				ImageName:   "nginx",
				Memory:      1.5,
				NumberCores: 1,
				Ports:       []int{80, 443},
			},
			bindingParameters: &aci.BindingParameters{},
		},
	}, nil
}

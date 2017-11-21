// +build !unit

package lifecycle

import (
	ac "github.com/Azure/azure-service-broker/pkg/azure/aci"
	"github.com/Azure/azure-service-broker/pkg/azure/arm"
	"github.com/Azure/azure-service-broker/pkg/services/aci"
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
			provisioningParameters: &aci.ProvisioningParameters{
				Location:    "eastus",
				ImageName:   "nginx",
				Memory:      1.5,
				NumberCores: 1,
				Ports:       []int{80, 443},
			},
			bindingParameters: &aci.BindingParameters{},
		},
	}, nil
}

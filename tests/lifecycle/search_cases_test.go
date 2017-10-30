// +build !unit

package lifecycle

import (
	"github.com/Azure/azure-service-broker/pkg/azure/arm"
	as "github.com/Azure/azure-service-broker/pkg/azure/search"
	"github.com/Azure/azure-service-broker/pkg/services/search"
)

func getSearchCases(
	armDeployer arm.Deployer,
	resourceGroup string,
) ([]moduleLifecycleTestCase, error) {
	searchManager, err := as.NewManager()
	if err != nil {
		return nil, err
	}

	return []moduleLifecycleTestCase{
		{
			module:    search.New(armDeployer, searchManager),
			serviceID: "c54902aa-3027-4c5c-8e96-5b3d3b452f7f",
			planID:    "35bd6e80-5ff5-487e-be0e-338aee9321e4",
			provisioningParameters: &search.ProvisioningParameters{
				Location: "southcentralus",
			},
			bindingParameters: &search.BindingParameters{},
		},
	}, nil
}

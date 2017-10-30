// +build !unit

package lifecycle

import (
	"github.com/Azure/azure-service-broker/pkg/azure/arm"
	rc "github.com/Azure/azure-service-broker/pkg/azure/rediscache"
	"github.com/Azure/azure-service-broker/pkg/services/rediscache"
)

func getRediscacheCases(
	armDeployer arm.Deployer,
	resourceGroup string,
) ([]moduleLifecycleTestCase, error) {
	redisManager, err := rc.NewManager()
	if err != nil {
		return nil, err
	}

	return []moduleLifecycleTestCase{
		{
			module:    rediscache.New(armDeployer, redisManager),
			serviceID: "0346088a-d4b2-4478-aa32-f18e295ec1d9",
			planID:    "362b3d1b-5b57-4289-80ad-4a15a760c29c",
			provisioningParameters: &rediscache.ProvisioningParameters{
				Location: "southcentralus",
			},
			bindingParameters: &rediscache.BindingParameters{},
		},
	}, nil
}

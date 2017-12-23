// +build !unit

package lifecycle

import (
	"github.com/Azure/open-service-broker-azure/pkg/azure/arm"
	rc "github.com/Azure/open-service-broker-azure/pkg/azure/rediscache"
	"github.com/Azure/open-service-broker-azure/pkg/services/rediscache"
)

func getRediscacheCases(
	armDeployer arm.Deployer,
	resourceGroup string,
) ([]serviceLifecycleTestCase, error) {
	redisManager, err := rc.NewManager()
	if err != nil {
		return nil, err
	}

	return []serviceLifecycleTestCase{
		{
			module:                 rediscache.New(armDeployer, redisManager),
			serviceID:              "0346088a-d4b2-4478-aa32-f18e295ec1d9",
			planID:                 "362b3d1b-5b57-4289-80ad-4a15a760c29c",
			location:               "southcentralus",
			provisioningParameters: &rediscache.ProvisioningParameters{},
			bindingParameters:      &rediscache.BindingParameters{},
		},
	}, nil
}

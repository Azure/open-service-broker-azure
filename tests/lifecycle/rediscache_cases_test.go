// +build !unit

package lifecycle

import (
	redisSDK "github.com/Azure/azure-sdk-for-go/services/redis/mgmt/2017-10-01/redis" // nolint: lll
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/open-service-broker-azure/pkg/azure/arm"
	"github.com/Azure/open-service-broker-azure/pkg/services/rediscache"
)

func getRediscacheCases(
	azureEnvironment azure.Environment,
	subscriptionID string,
	authorizer autorest.Authorizer,
	armDeployer arm.Deployer,
) ([]serviceLifecycleTestCase, error) {
	client := redisSDK.NewClientWithBaseURI(
		azureEnvironment.ResourceManagerEndpoint,
		subscriptionID,
	)
	client.Authorizer = authorizer
	return []serviceLifecycleTestCase{
		{
			module:                 rediscache.New(armDeployer, client),
			serviceID:              "0346088a-d4b2-4478-aa32-f18e295ec1d9",
			planID:                 "362b3d1b-5b57-4289-80ad-4a15a760c29c",
			location:               "southcentralus",
			provisioningParameters: &rediscache.ProvisioningParameters{},
			bindingParameters:      &rediscache.BindingParameters{},
		},
	}, nil
}

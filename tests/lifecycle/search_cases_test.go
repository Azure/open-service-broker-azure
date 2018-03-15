// +build !unit

package lifecycle

import (
	searchSDK "github.com/Azure/azure-sdk-for-go/services/search/mgmt/2015-08-19/search" // nolint: lll
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/open-service-broker-azure/pkg/azure/arm"
	"github.com/Azure/open-service-broker-azure/pkg/services/search"
)

func getSearchCases(
	azureEnvironment azure.Environment,
	subscriptionID string,
	authorizer autorest.Authorizer,
	armDeployer arm.Deployer,
) ([]serviceLifecycleTestCase, error) {
	servicesClient := searchSDK.NewServicesClientWithBaseURI(
		azureEnvironment.ResourceManagerEndpoint,
		subscriptionID,
	)
	servicesClient.Authorizer = authorizer
	return []serviceLifecycleTestCase{
		{
			module:    search.New(armDeployer, servicesClient),
			serviceID: "c54902aa-3027-4c5c-8e96-5b3d3b452f7f",
			planID:    "35bd6e80-5ff5-487e-be0e-338aee9321e4",
			location:  "southcentralus",
		},
	}, nil
}

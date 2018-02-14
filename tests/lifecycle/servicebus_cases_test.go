// +build !unit

package lifecycle

import (
	servicebusSDK "github.com/Azure/azure-sdk-for-go/arm/servicebus"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/open-service-broker-azure/pkg/azure/arm"
	"github.com/Azure/open-service-broker-azure/pkg/services/servicebus"
)

func getServicebusCases(
	azureEnvironment azure.Environment,
	subscriptionID string,
	authorizer autorest.Authorizer,
	armDeployer arm.Deployer,
) ([]serviceLifecycleTestCase, error) {
	namespacesClient := servicebusSDK.NewNamespacesClientWithBaseURI(
		azureEnvironment.ResourceManagerEndpoint,
		subscriptionID,
	)
	namespacesClient.Authorizer = authorizer
	return []serviceLifecycleTestCase{
		{
			module:                 servicebus.New(armDeployer, namespacesClient),
			serviceID:              "6dc44338-2f13-4bc5-9247-5b1b3c5462d3",
			planID:                 "d06817b1-87ea-4320-8942-14b1d060206a",
			location:               "southcentralus",
			provisioningParameters: &servicebus.ProvisioningParameters{},
			bindingParameters:      &servicebus.BindingParameters{},
		},
	}, nil
}

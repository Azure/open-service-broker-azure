// +build !unit

package lifecycle

import (
	eventHubSDK "github.com/Azure/azure-sdk-for-go/services/eventhub/mgmt/2017-04-01/eventhub" // nolint: lll
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/open-service-broker-azure/pkg/azure/arm"
	"github.com/Azure/open-service-broker-azure/pkg/services/eventhubs"
)

func getEventhubCases(
	azureEnvironment azure.Environment,
	subscriptionID string,
	authorizer autorest.Authorizer,
	armDeployer arm.Deployer,
) ([]serviceLifecycleTestCase, error) {
	namespacesClient := eventHubSDK.NewNamespacesClientWithBaseURI(
		azureEnvironment.ResourceManagerEndpoint,
		subscriptionID,
	)
	namespacesClient.Authorizer = authorizer
	return []serviceLifecycleTestCase{
		{
			module:    eventhubs.New(armDeployer, namespacesClient),
			serviceID: "7bade660-32f1-4fd7-b9e6-d416d975170b",
			planID:    "80756db5-a20c-495d-ae70-62cf7d196a3c",
			location:  "southcentralus",
		},
	}, nil
}

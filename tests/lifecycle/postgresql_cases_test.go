// +build !unit

package lifecycle

import (
	postgresSDK "github.com/Azure/azure-sdk-for-go/services/postgresql/mgmt/2017-04-30-preview/postgresql" // nolint: lll
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/open-service-broker-azure/pkg/azure/arm"
	"github.com/Azure/open-service-broker-azure/pkg/services/postgresqldb"
)

func getPostgresqlCases(
	azureEnvironment azure.Environment,
	subscriptionID string,
	authorizer autorest.Authorizer,
	armDeployer arm.Deployer,
) ([]serviceLifecycleTestCase, error) {
	checkNameAvailabilityClient :=
		postgresSDK.NewCheckNameAvailabilityClientWithBaseURI(
			azureEnvironment.ResourceManagerEndpoint,
			subscriptionID,
		)
	checkNameAvailabilityClient.Authorizer = authorizer
	serversClient := postgresSDK.NewServersClientWithBaseURI(
		azureEnvironment.ResourceManagerEndpoint,
		subscriptionID,
	)
	serversClient.Authorizer = authorizer
	return []serviceLifecycleTestCase{
		{
			module: postgresqldb.New(
				armDeployer,
				checkNameAvailabilityClient,
				serversClient,
			),
			serviceID: "b43b4bba-5741-4d98-a10b-17dc5cee0175",
			planID:    "b2ed210f-6a10-4593-a6c4-964e6b6fad62",
			location:  "southcentralus",
			provisioningParameters: &postgresqldb.ProvisioningParameters{
				FirewallIPStart: "0.0.0.0",
				FirewallIPEnd:   "255.255.255.255",
				Extensions: []string{
					"uuid-ossp",
					"postgis",
				},
			},
			bindingParameters: &postgresqldb.BindingParameters{},
		},
	}, nil
}

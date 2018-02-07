// +build !unit

package lifecycle

import (
	mysqlSDK "github.com/Azure/azure-sdk-for-go/arm/mysql"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/open-service-broker-azure/pkg/azure/arm"
	"github.com/Azure/open-service-broker-azure/pkg/services/mysqldb"
)

func getMysqlCases(
	azureEnvironment azure.Environment,
	subscriptionID string,
	authorizer autorest.Authorizer,
	armDeployer arm.Deployer,
) ([]serviceLifecycleTestCase, error) {
	serversClient := mysqlSDK.NewServersClientWithBaseURI(
		azureEnvironment.ResourceManagerEndpoint,
		subscriptionID,
	)
	serversClient.Authorizer = authorizer
	return []serviceLifecycleTestCase{
		{
			module:    mysqldb.New(armDeployer, serversClient),
			serviceID: "997b8372-8dac-40ac-ae65-758b4a5075a5",
			planID:    "427559f1-bf2a-45d3-8844-32374a3e58aa",
			location:  "southcentralus",
			provisioningParameters: &mysqldb.ProvisioningParameters{
				FirewallIPStart: "0.0.0.0",
				FirewallIPEnd:   "255.255.255.255",
			},
			bindingParameters: &mysqldb.BindingParameters{},
		},
	}, nil
}

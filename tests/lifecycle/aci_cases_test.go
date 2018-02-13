// +build !unit

package lifecycle

import (
	aciSDK "github.com/Azure/azure-sdk-for-go/services/containerinstance/mgmt/2017-08-01-preview/containerinstance" // nolint: lll
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/open-service-broker-azure/pkg/azure/arm"
	"github.com/Azure/open-service-broker-azure/pkg/services/aci"
)

func getACICases(
	azureEnvironment azure.Environment,
	subscriptionID string,
	authorizer autorest.Authorizer,
	armDeployer arm.Deployer,
) ([]serviceLifecycleTestCase, error) {
	containerGroupsClient := aciSDK.NewContainerGroupsClientWithBaseURI(
		azureEnvironment.ResourceManagerEndpoint,
		subscriptionID,
	)
	containerGroupsClient.Authorizer = authorizer
	return []serviceLifecycleTestCase{
		{
			module:    aci.New(armDeployer, containerGroupsClient),
			serviceID: "451d5d19-4575-4d4a-9474-116f705ecc95",
			planID:    "d48798e2-21db-405b-abc7-aa6f0ff08f6c",
			location:  "eastus",
			provisioningParameters: &aci.ProvisioningParameters{
				ImageName:   "nginx",
				Memory:      1.5,
				NumberCores: 1,
				Ports:       []int{80, 443},
			},
			bindingParameters: &aci.BindingParameters{},
		},
	}, nil
}

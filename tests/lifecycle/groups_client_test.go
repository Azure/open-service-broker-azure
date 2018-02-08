// +build !unit

package lifecycle

import (
	"github.com/Azure/azure-sdk-for-go/arm/resources/resources"
	az "github.com/Azure/open-service-broker-azure/pkg/azure"
	"github.com/Azure/open-service-broker-azure/pkg/config"
)

func ensureResourceGroup(resourceGroup string) error {
	groupsClient, err := getGroupsClient()
	if err != nil {
		return err
	}
	location := "eastus"
	_, err = groupsClient.CreateOrUpdate(
		resourceGroup,
		resources.Group{
			Name:     &resourceGroup,
			Location: &location,
		},
	)
	return err
}

func deleteResourceGroup(
	resourceGroupName string,
) error {
	groupsClient, err := getGroupsClient()
	if err != nil {
		return err
	}
	cancelCh := make(chan struct{})
	defer close(cancelCh)
	_, errCh := groupsClient.Delete(resourceGroupName, cancelCh)
	return <-errCh
}

func getGroupsClient() (*resources.GroupsClient, error) {
	azureConfig, err := config.GetAzureConfig()
	if err != nil {
		return nil, err
	}
	groupsClient := resources.NewGroupsClientWithBaseURI(
		azureConfig.Environment.ResourceManagerEndpoint,
		azureConfig.SubscriptionID,
	)
	authorizer, err := az.GetBearerTokenAuthorizer(
		azureConfig.Environment,
		azureConfig.TenantID,
		azureConfig.ClientID,
		azureConfig.ClientSecret,
	)
	if err != nil {
		return nil, err
	}
	groupsClient.Authorizer = authorizer
	return &groupsClient, err
}

// +build !unit

package lifecycle

import (
	"github.com/Azure/azure-sdk-for-go/arm/resources/resources"
	az "github.com/Azure/open-service-broker-azure/pkg/azure"
	"github.com/Azure/go-autorest/autorest/azure"
)

func ensureResourceGroup(resourceGroup string, location string) error {
	groupsClient, err := getGroupsClient()
	if err != nil {
		return err
	}
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
	azureConfig, err := az.GetConfig()
	if err != nil {
		return nil, err
	}
	azureEnvironment, err := azure.EnvironmentFromName(azureConfig.Environment)
	if err != nil {
		return nil, err
	}
	groupsClient := resources.NewGroupsClientWithBaseURI(
		azureEnvironment.ResourceManagerEndpoint,
		azureConfig.SubscriptionID,
	)
	authorizer, err := az.GetBearerTokenAuthorizer(
		azureEnvironment,
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

// +build !unit

package lifecycle

import (
	"context"

	resourcesSDK "github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2017-05-10/resources" // nolint: lll
	"open-service-broker-azure/pkg/azure"
)

func ensureResourceGroup(resourceGroup string) error {
	groupsClient, err := getGroupsClient()
	if err != nil {
		return err
	}
	location := "eastus"
	_, err = groupsClient.CreateOrUpdate(
		context.Background(),
		resourceGroup,
		resourcesSDK.Group{
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
	_, err = groupsClient.Delete(context.Background(), resourceGroupName)
	return err
}

func getGroupsClient() (*resourcesSDK.GroupsClient, error) {
	azureConfig, err := azure.GetConfigFromEnvironment()
	if err != nil {
		return nil, err
	}
	groupsClient := resourcesSDK.NewGroupsClientWithBaseURI(
		azureConfig.Environment.ResourceManagerEndpoint,
		azureConfig.SubscriptionID,
	)
	authorizer, err := azure.GetBearerTokenAuthorizer(
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

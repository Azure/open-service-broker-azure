// +build !unit

package lifecycle

import (
	"context"

	resourcesSDK "github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2017-05-10/resources" // nolint: lll
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
	azureConfig, err := config.GetAzureConfig()
	if err != nil {
		return nil, err
	}
	groupsClient := resourcesSDK.NewGroupsClientWithBaseURI(
		azureConfig.GetEnvironment().ResourceManagerEndpoint,
		azureConfig.GetSubscriptionID(),
	)
	authorizer, err := az.GetBearerTokenAuthorizer(
		azureConfig.GetEnvironment(),
		azureConfig.GetTenantID(),
		azureConfig.GetClientID(),
		azureConfig.GetClientSecret(),
	)
	if err != nil {
		return nil, err
	}
	groupsClient.Authorizer = authorizer
	return &groupsClient, err
}

// +build !unit

package lifecycle

import (
	"github.com/Azure/azure-sdk-for-go/arm/resources/resources"
	az "github.com/Azure/azure-service-broker/pkg/azure"
	"github.com/Azure/go-autorest/autorest/azure"
)

func deleteResourceGroup(resourceGroupName string) error {
	azureConfig, err := az.GetConfig()
	if err != nil {
		return err
	}
	azureEnvironment, err := azure.EnvironmentFromName(azureConfig.Environment)
	if err != nil {
		return err
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
		return err
	}
	groupsClient.Authorizer = authorizer
	cancelCh := make(chan struct{})
	defer close(cancelCh)
	_, errChan := groupsClient.Delete(resourceGroupName, cancelCh)
	return <-errChan
}

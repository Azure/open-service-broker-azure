package boot

import (
	"fmt"
	"time"

	postgresSDK "github.com/Azure/azure-sdk-for-go/services/postgresql/mgmt/2017-04-30-preview/postgresql" // nolint: lll
	resourcesSDK "github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2017-05-10/resources"          // nolint: lll
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/open-service-broker-azure/pkg/azure"
	"github.com/Azure/open-service-broker-azure/pkg/azure/arm"
	"github.com/Azure/open-service-broker-azure/pkg/service"
	"github.com/Azure/open-service-broker-azure/pkg/services/postgresql"
	"github.com/Azure/open-service-broker-azure/pkg/version"
)

func getModules(
	catalogConfig service.CatalogConfig,
	azureConfig azure.Config,
) ([]service.Module, error) {
	azureSubscriptionID := azureConfig.SubscriptionID

	authorizer, err := azure.GetBearerTokenAuthorizer(
		azureConfig.Environment,
		azureConfig.TenantID,
		azureConfig.ClientID,
		azureConfig.ClientSecret,
	)
	if err != nil {
		return nil, fmt.Errorf("error getting bearer token authorizer: %s", err)
	}

	resourceGroupsClient := resourcesSDK.NewGroupsClientWithBaseURI(
		azureConfig.Environment.ResourceManagerEndpoint,
		azureSubscriptionID,
	)
	resourceGroupsClient.Authorizer = authorizer
	resourceGroupsClient.UserAgent = getUserAgent(resourceGroupsClient.Client)
	resourceDeploymentsClient := resourcesSDK.NewDeploymentsClientWithBaseURI(
		azureConfig.Environment.ResourceManagerEndpoint,
		azureSubscriptionID,
	)
	resourceDeploymentsClient.Authorizer = authorizer
	resourceDeploymentsClient.UserAgent =
		getUserAgent(resourceDeploymentsClient.Client)
	resourceDeploymentsClient.PollingDuration = time.Minute * 30
	armDeployer := arm.NewDeployer(
		resourceGroupsClient,
		resourceDeploymentsClient,
	)

	postgresCheckNameAvailabilityClient :=
		postgresSDK.NewCheckNameAvailabilityClientWithBaseURI(
			azureConfig.Environment.ResourceManagerEndpoint,
			azureSubscriptionID,
		)
	postgresCheckNameAvailabilityClient.Authorizer = authorizer
	postgresCheckNameAvailabilityClient.UserAgent =
		getUserAgent(postgresCheckNameAvailabilityClient.Client)
	postgresServersClient := postgresSDK.NewServersClientWithBaseURI(
		azureConfig.Environment.ResourceManagerEndpoint,
		azureSubscriptionID,
	)
	postgresServersClient.Authorizer = authorizer
	postgresServersClient.UserAgent = getUserAgent(postgresServersClient.Client)
	postgresDatabasesClient := postgresSDK.NewDatabasesClientWithBaseURI(
		azureConfig.Environment.ResourceManagerEndpoint,
		azureSubscriptionID,
	)
	postgresDatabasesClient.Authorizer = authorizer
	postgresDatabasesClient.UserAgent = getUserAgent(postgresServersClient.Client)

	modules := []service.Module{
		postgresql.New(
			armDeployer,
			postgresCheckNameAvailabilityClient,
			postgresServersClient,
			postgresDatabasesClient,
		),
	}

	// Filter modules based on stability
	filteredModules := []service.Module{}
	for _, module := range modules {
		if module.GetStability() >= catalogConfig.MinStability {
			filteredModules = append(filteredModules, module)
		}
	}

	return filteredModules, nil
}

func getUserAgent(client autorest.Client) string {
	return fmt.Sprintf(
		"%s; open-service-broker/%s",
		client.UserAgent,
		version.GetVersion(),
	)
}

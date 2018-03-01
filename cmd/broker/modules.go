package main

import (
	"fmt"
	"time"

	aciSDK "github.com/Azure/azure-sdk-for-go/services/containerinstance/mgmt/2017-08-01-preview/containerinstance" // nolint: lll
	cosmosSDK "github.com/Azure/azure-sdk-for-go/services/cosmos-db/mgmt/2015-04-08/documentdb"                     // nolint: lll
	eventHubSDK "github.com/Azure/azure-sdk-for-go/services/eventhub/mgmt/2017-04-01/eventhub"                      // nolint: lll
	keyVaultSDK "github.com/Azure/azure-sdk-for-go/services/keyvault/mgmt/2016-10-01/keyvault"                      // nolint: lll
	mysqlSDK "github.com/Azure/azure-sdk-for-go/services/mysql/mgmt/2017-04-30-preview/mysql"                       // nolint: lll
	postgresSDK "github.com/Azure/azure-sdk-for-go/services/postgresql/mgmt/2017-04-30-preview/postgresql"          // nolint: lll
	redisSDK "github.com/Azure/azure-sdk-for-go/services/redis/mgmt/2017-10-01/redis"                               // nolint: lll
	resourcesSDK "github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2017-05-10/resources"                   // nolint: lll
	searchSDK "github.com/Azure/azure-sdk-for-go/services/search/mgmt/2015-08-19/search"                            // nolint: lll
	servicebusSDK "github.com/Azure/azure-sdk-for-go/services/servicebus/mgmt/2017-04-01/servicebus"                // nolint: lll
	sqlSDK "github.com/Azure/azure-sdk-for-go/services/sql/mgmt/2017-03-01-preview/sql"                             // nolint: lll
	storageSDK "github.com/Azure/azure-sdk-for-go/services/storage/mgmt/2017-10-01/storage"                         // nolint: lll
	"github.com/Azure/go-autorest/autorest"
	az "github.com/Azure/open-service-broker-azure/pkg/azure"
	"github.com/Azure/open-service-broker-azure/pkg/azure/arm"
	"github.com/Azure/open-service-broker-azure/pkg/config"
	"github.com/Azure/open-service-broker-azure/pkg/service"
	"github.com/Azure/open-service-broker-azure/pkg/services/aci"
	"github.com/Azure/open-service-broker-azure/pkg/services/cosmosdb"
	"github.com/Azure/open-service-broker-azure/pkg/services/eventhubs"
	"github.com/Azure/open-service-broker-azure/pkg/services/keyvault"
	"github.com/Azure/open-service-broker-azure/pkg/services/mssql"
	"github.com/Azure/open-service-broker-azure/pkg/services/mysql"
	"github.com/Azure/open-service-broker-azure/pkg/services/postgresql"
	"github.com/Azure/open-service-broker-azure/pkg/services/rediscache"
	"github.com/Azure/open-service-broker-azure/pkg/services/search"
	"github.com/Azure/open-service-broker-azure/pkg/services/servicebus"
	"github.com/Azure/open-service-broker-azure/pkg/services/storage"
	"github.com/Azure/open-service-broker-azure/pkg/version"
)

var modules []service.Module

func initModules(azureConfig config.AzureConfig) error {
	azureEnvironment := azureConfig.GetEnvironment()
	azureSubscriptionID := azureConfig.GetSubscriptionID()

	authorizer, err := az.GetBearerTokenAuthorizer(
		azureEnvironment,
		azureConfig.GetTenantID(),
		azureConfig.GetClientID(),
		azureConfig.GetClientSecret(),
	)
	if err != nil {
		return fmt.Errorf("error getting bearer token authorizer: %s", err)
	}

	resourceGroupsClient := resourcesSDK.NewGroupsClientWithBaseURI(
		azureEnvironment.ResourceManagerEndpoint,
		azureSubscriptionID,
	)
	resourceGroupsClient.Authorizer = authorizer
	resourceGroupsClient.UserAgent = getUserAgent(resourceGroupsClient.Client)
	resourceDeploymentsClient := resourcesSDK.NewDeploymentsClientWithBaseURI(
		azureEnvironment.ResourceManagerEndpoint,
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

	aciClient := aciSDK.NewContainerGroupsClientWithBaseURI(
		azureEnvironment.ResourceManagerEndpoint,
		azureSubscriptionID,
	)
	aciClient.Authorizer = authorizer
	aciClient.UserAgent = getUserAgent(aciClient.Client)

	cosmosdbAccountsClient := cosmosSDK.NewDatabaseAccountsClientWithBaseURI(
		azureEnvironment.ResourceManagerEndpoint,
		azureSubscriptionID,
	)
	cosmosdbAccountsClient.Authorizer = authorizer
	cosmosdbAccountsClient.UserAgent = getUserAgent(cosmosdbAccountsClient.Client)

	eventHubNamespacesClient := eventHubSDK.NewNamespacesClientWithBaseURI(
		azureEnvironment.ResourceManagerEndpoint,
		azureSubscriptionID,
	)
	eventHubNamespacesClient.Authorizer = authorizer
	eventHubNamespacesClient.UserAgent =
		getUserAgent(eventHubNamespacesClient.Client)

	keyVaultsClient := keyVaultSDK.NewVaultsClientWithBaseURI(
		azureEnvironment.ResourceManagerEndpoint,
		azureSubscriptionID,
	)
	keyVaultsClient.Authorizer = authorizer
	keyVaultsClient.UserAgent = getUserAgent(keyVaultsClient.Client)

	mysqlCheckNameAvailabilityClient :=
		mysqlSDK.NewCheckNameAvailabilityClientWithBaseURI(
			azureEnvironment.ResourceManagerEndpoint,
			azureSubscriptionID,
		)
	mysqlCheckNameAvailabilityClient.Authorizer = authorizer
	mysqlCheckNameAvailabilityClient.UserAgent =
		getUserAgent(mysqlCheckNameAvailabilityClient.Client)
	mysqlServersClient := mysqlSDK.NewServersClientWithBaseURI(
		azureEnvironment.ResourceManagerEndpoint,
		azureSubscriptionID,
	)
	mysqlServersClient.Authorizer = authorizer
	mysqlServersClient.UserAgent = getUserAgent(mysqlServersClient.Client)
	mysqlDatabasesClient := mysqlSDK.NewDatabasesClientWithBaseURI(
		azureEnvironment.ResourceManagerEndpoint,
		azureSubscriptionID,
	)
	mysqlDatabasesClient.Authorizer = authorizer
	mysqlDatabasesClient.UserAgent = getUserAgent(mysqlDatabasesClient.Client)

	postgresCheckNameAvailabilityClient :=
		postgresSDK.NewCheckNameAvailabilityClientWithBaseURI(
			azureEnvironment.ResourceManagerEndpoint,
			azureSubscriptionID,
		)
	postgresCheckNameAvailabilityClient.Authorizer = authorizer
	postgresCheckNameAvailabilityClient.UserAgent =
		getUserAgent(postgresCheckNameAvailabilityClient.Client)
	postgresServersClient := postgresSDK.NewServersClientWithBaseURI(
		azureEnvironment.ResourceManagerEndpoint,
		azureSubscriptionID,
	)
	postgresServersClient.Authorizer = authorizer
	postgresServersClient.UserAgent = getUserAgent(postgresServersClient.Client)
	postgresDatabasesClient := postgresSDK.NewDatabasesClientWithBaseURI(
		azureEnvironment.ResourceManagerEndpoint,
		azureSubscriptionID,
	)
	postgresDatabasesClient.Authorizer = authorizer
	postgresDatabasesClient.UserAgent = getUserAgent(postgresServersClient.Client)

	sqlServersClient := sqlSDK.NewServersClientWithBaseURI(
		azureEnvironment.ResourceManagerEndpoint,
		azureSubscriptionID,
	)
	sqlServersClient.Authorizer = authorizer
	sqlServersClient.UserAgent = getUserAgent(sqlServersClient.Client)
	sqlDatabasesClient := sqlSDK.NewDatabasesClientWithBaseURI(
		azureEnvironment.ResourceManagerEndpoint,
		azureSubscriptionID,
	)
	sqlDatabasesClient.Authorizer = authorizer
	sqlDatabasesClient.UserAgent = getUserAgent(sqlDatabasesClient.Client)

	redisClient := redisSDK.NewClientWithBaseURI(
		azureEnvironment.ResourceManagerEndpoint,
		azureSubscriptionID,
	)
	redisClient.Authorizer = authorizer
	redisClient.UserAgent = getUserAgent(redisClient.Client)

	searchServicesClient := searchSDK.NewServicesClientWithBaseURI(
		azureEnvironment.ResourceManagerEndpoint,
		azureSubscriptionID,
	)
	searchServicesClient.Authorizer = authorizer
	searchServicesClient.UserAgent = getUserAgent(searchServicesClient.Client)

	serviceBusNamespacesClient := servicebusSDK.NewNamespacesClientWithBaseURI(
		azureEnvironment.ResourceManagerEndpoint,
		azureSubscriptionID,
	)
	serviceBusNamespacesClient.Authorizer = authorizer
	serviceBusNamespacesClient.UserAgent =
		getUserAgent(serviceBusNamespacesClient.Client)

	storageAccountsClient := storageSDK.NewAccountsClientWithBaseURI(
		azureEnvironment.ResourceManagerEndpoint,
		azureSubscriptionID,
	)
	storageAccountsClient.Authorizer = authorizer
	storageAccountsClient.UserAgent = getUserAgent(storageAccountsClient.Client)

	modules = []service.Module{
		postgresql.New(
			armDeployer,
			postgresCheckNameAvailabilityClient,
			postgresServersClient,
			postgresDatabasesClient,
		),
		rediscache.New(armDeployer, redisClient),
		mysql.New(
			azureEnvironment,
			armDeployer,
			mysqlCheckNameAvailabilityClient,
			mysqlServersClient,
			mysqlDatabasesClient,
		),
		servicebus.New(armDeployer, serviceBusNamespacesClient),
		eventhubs.New(armDeployer, eventHubNamespacesClient),
		keyvault.New(azureConfig.GetTenantID(), armDeployer, keyVaultsClient),
		mssql.New(
			azureEnvironment,
			armDeployer,
			sqlServersClient,
			sqlDatabasesClient,
		),
		cosmosdb.New(armDeployer, cosmosdbAccountsClient),
		storage.New(armDeployer, storageAccountsClient),
		search.New(armDeployer, searchServicesClient),
		aci.New(armDeployer, aciClient),
	}
	return nil
}

func getUserAgent(client autorest.Client) string {
	return fmt.Sprintf(
		"%s; open-service-broker/%s",
		client.UserAgent,
		version.GetVersion(),
	)
}

package main

import (
	"fmt"

	aciSDK "github.com/Azure/azure-sdk-for-go/arm/containerinstance"
	cosmosSDK "github.com/Azure/azure-sdk-for-go/arm/cosmos-db"
	eventHubSDK "github.com/Azure/azure-sdk-for-go/arm/eventhub"
	keyVaultSDK "github.com/Azure/azure-sdk-for-go/arm/keyvault"
	mysqlSDK "github.com/Azure/azure-sdk-for-go/arm/mysql"
	postgresSDK "github.com/Azure/azure-sdk-for-go/arm/postgresql"
	redisSDK "github.com/Azure/azure-sdk-for-go/arm/redis"
	"github.com/Azure/azure-sdk-for-go/arm/resources/resources"
	searchSDK "github.com/Azure/azure-sdk-for-go/arm/search"
	servicebusSDK "github.com/Azure/azure-sdk-for-go/arm/servicebus"
	sqlSDK "github.com/Azure/azure-sdk-for-go/arm/sql"
	storageSDK "github.com/Azure/azure-sdk-for-go/arm/storage"
	"github.com/Azure/go-autorest/autorest"
	az "github.com/Azure/open-service-broker-azure/pkg/azure"
	"github.com/Azure/open-service-broker-azure/pkg/azure/arm"
	"github.com/Azure/open-service-broker-azure/pkg/config"
	"github.com/Azure/open-service-broker-azure/pkg/service"
	"github.com/Azure/open-service-broker-azure/pkg/services/aci"
	"github.com/Azure/open-service-broker-azure/pkg/services/cosmosdb"
	"github.com/Azure/open-service-broker-azure/pkg/services/eventhubs"
	"github.com/Azure/open-service-broker-azure/pkg/services/keyvault"
	"github.com/Azure/open-service-broker-azure/pkg/services/mysqldb"
	"github.com/Azure/open-service-broker-azure/pkg/services/postgresqldb"
	"github.com/Azure/open-service-broker-azure/pkg/services/rediscache"
	"github.com/Azure/open-service-broker-azure/pkg/services/search"
	"github.com/Azure/open-service-broker-azure/pkg/services/servicebus"
	"github.com/Azure/open-service-broker-azure/pkg/services/sqldb"
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

	resourceGroupsClient := resources.NewGroupsClientWithBaseURI(
		azureEnvironment.ResourceManagerEndpoint,
		azureSubscriptionID,
	)
	resourceGroupsClient.Authorizer = authorizer
	resourceGroupsClient.UserAgent = getUserAgent(resourceGroupsClient.Client)
	resourceDeploymentsClient := resources.NewDeploymentsClientWithBaseURI(
		azureEnvironment.ResourceManagerEndpoint,
		azureSubscriptionID,
	)
	resourceDeploymentsClient.Authorizer = authorizer
	resourceDeploymentsClient.UserAgent =
		getUserAgent(resourceDeploymentsClient.Client)
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

	mysqlServersClient := mysqlSDK.NewServersClientWithBaseURI(
		azureEnvironment.ResourceManagerEndpoint,
		azureSubscriptionID,
	)
	mysqlServersClient.Authorizer = authorizer
	mysqlServersClient.UserAgent = getUserAgent(mysqlServersClient.Client)

	postgresServersClient := postgresSDK.NewServersClientWithBaseURI(
		azureEnvironment.ResourceManagerEndpoint,
		azureSubscriptionID,
	)
	postgresServersClient.Authorizer = authorizer
	postgresServersClient.UserAgent = getUserAgent(postgresServersClient.Client)

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

	redisGroupClient := redisSDK.NewGroupClientWithBaseURI(
		azureEnvironment.ResourceManagerEndpoint,
		azureSubscriptionID,
	)
	redisGroupClient.Authorizer = authorizer
	redisGroupClient.UserAgent = getUserAgent(redisGroupClient.Client)

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
		postgresqldb.New(armDeployer, postgresServersClient),
		rediscache.New(armDeployer, redisGroupClient),
		mysqldb.New(azureEnvironment, armDeployer, mysqlServersClient),
		servicebus.New(armDeployer, serviceBusNamespacesClient),
		eventhubs.New(armDeployer, eventHubNamespacesClient),
		keyvault.New(azureConfig.GetTenantID(), armDeployer, keyVaultsClient),
		sqldb.New(
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

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
	authorizer, err := az.GetBearerTokenAuthorizer(
		azureConfig.Environment,
		azureConfig.TenantID,
		azureConfig.ClientID,
		azureConfig.ClientSecret,
	)
	if err != nil {
		return fmt.Errorf("error getting bearer token authorizer: %s", err)
	}

	resourceGroupsClient := resources.NewGroupsClientWithBaseURI(
		azureConfig.Environment.ResourceManagerEndpoint,
		azureConfig.SubscriptionID,
	)
	resourceGroupsClient.Authorizer = authorizer
	resourceGroupsClient.UserAgent = getUserAgent(resourceGroupsClient.Client)
	resourceDeploymentsClient := resources.NewDeploymentsClientWithBaseURI(
		azureConfig.Environment.ResourceManagerEndpoint,
		azureConfig.SubscriptionID,
	)
	resourceDeploymentsClient.Authorizer = authorizer
	resourceDeploymentsClient.UserAgent =
		getUserAgent(resourceDeploymentsClient.Client)
	armDeployer := arm.NewDeployer(
		resourceGroupsClient,
		resourceDeploymentsClient,
	)

	aciClient := aciSDK.NewContainerGroupsClientWithBaseURI(
		azureConfig.Environment.ResourceManagerEndpoint,
		azureConfig.SubscriptionID,
	)
	aciClient.Authorizer = authorizer
	aciClient.UserAgent = getUserAgent(aciClient.Client)

	cosmosdbAccountsClient := cosmosSDK.NewDatabaseAccountsClientWithBaseURI(
		azureConfig.Environment.ResourceManagerEndpoint,
		azureConfig.SubscriptionID,
	)
	cosmosdbAccountsClient.Authorizer = authorizer
	cosmosdbAccountsClient.UserAgent = getUserAgent(cosmosdbAccountsClient.Client)

	eventHubNamespacesClient := eventHubSDK.NewNamespacesClientWithBaseURI(
		azureConfig.Environment.ResourceManagerEndpoint,
		azureConfig.SubscriptionID,
	)
	eventHubNamespacesClient.Authorizer = authorizer
	eventHubNamespacesClient.UserAgent =
		getUserAgent(eventHubNamespacesClient.Client)

	keyVaultsClient := keyVaultSDK.NewVaultsClientWithBaseURI(
		azureConfig.Environment.ResourceManagerEndpoint,
		azureConfig.SubscriptionID,
	)
	keyVaultsClient.Authorizer = authorizer
	keyVaultsClient.UserAgent = getUserAgent(keyVaultsClient.Client)

	mysqlServersClient := mysqlSDK.NewServersClientWithBaseURI(
		azureConfig.Environment.ResourceManagerEndpoint,
		azureConfig.SubscriptionID,
	)
	mysqlServersClient.Authorizer = authorizer
	mysqlServersClient.UserAgent = getUserAgent(mysqlServersClient.Client)

	postgresServersClient := postgresSDK.NewServersClientWithBaseURI(
		azureConfig.Environment.ResourceManagerEndpoint,
		azureConfig.SubscriptionID,
	)
	postgresServersClient.Authorizer = authorizer
	postgresServersClient.UserAgent = getUserAgent(postgresServersClient.Client)

	sqlServersClient := sqlSDK.NewServersClientWithBaseURI(
		azureConfig.Environment.ResourceManagerEndpoint,
		azureConfig.SubscriptionID,
	)
	sqlServersClient.Authorizer = authorizer
	sqlServersClient.UserAgent = getUserAgent(sqlServersClient.Client)
	sqlDatabasesClient := sqlSDK.NewDatabasesClientWithBaseURI(
		azureConfig.Environment.ResourceManagerEndpoint,
		azureConfig.SubscriptionID,
	)
	sqlDatabasesClient.Authorizer = authorizer
	sqlDatabasesClient.UserAgent = getUserAgent(sqlDatabasesClient.Client)

	redisGroupClient := redisSDK.NewGroupClientWithBaseURI(
		azureConfig.Environment.ResourceManagerEndpoint,
		azureConfig.SubscriptionID,
	)
	redisGroupClient.Authorizer = authorizer
	redisGroupClient.UserAgent = getUserAgent(redisGroupClient.Client)

	searchServicesClient := searchSDK.NewServicesClientWithBaseURI(
		azureConfig.Environment.ResourceManagerEndpoint,
		azureConfig.SubscriptionID,
	)
	searchServicesClient.Authorizer = authorizer
	searchServicesClient.UserAgent = getUserAgent(searchServicesClient.Client)

	serviceBusNamespacesClient := servicebusSDK.NewNamespacesClientWithBaseURI(
		azureConfig.Environment.ResourceManagerEndpoint,
		azureConfig.SubscriptionID,
	)
	serviceBusNamespacesClient.Authorizer = authorizer
	serviceBusNamespacesClient.UserAgent =
		getUserAgent(serviceBusNamespacesClient.Client)

	storageAccountsClient := storageSDK.NewAccountsClientWithBaseURI(
		azureConfig.Environment.ResourceManagerEndpoint,
		azureConfig.SubscriptionID,
	)
	storageAccountsClient.Authorizer = authorizer
	storageAccountsClient.UserAgent = getUserAgent(storageAccountsClient.Client)

	modules = []service.Module{
		postgresqldb.New(armDeployer, postgresServersClient),
		rediscache.New(armDeployer, redisGroupClient),
		mysqldb.New(azureConfig.Environment, armDeployer, mysqlServersClient),
		servicebus.New(armDeployer, serviceBusNamespacesClient),
		eventhubs.New(armDeployer, eventHubNamespacesClient),
		keyvault.New(azureConfig.TenantID, armDeployer, keyVaultsClient),
		sqldb.New(
			azureConfig.Environment,
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

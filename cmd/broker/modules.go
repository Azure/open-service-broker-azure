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
	searchSDK "github.com/Azure/azure-sdk-for-go/arm/search"
	servicebusSDK "github.com/Azure/azure-sdk-for-go/arm/servicebus"
	sqlSDK "github.com/Azure/azure-sdk-for-go/arm/sql"
	storageSDK "github.com/Azure/azure-sdk-for-go/arm/storage"
	"github.com/Azure/go-autorest/autorest/azure"
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
)

var modules []service.Module

func initModules() error {
	azureConfig, err := config.GetAzureConfig()
	if err != nil {
		return fmt.Errorf("error getting azure configuration: %s", err)
	}

	azureEnvironment, err := azure.EnvironmentFromName(azureConfig.Environment)
	if err != nil {
		return fmt.Errorf(
			"error getting azure environment from environment name: %s",
			err,
		)
	}

	authorizer, err := az.GetBearerTokenAuthorizer(
		azureEnvironment,
		azureConfig.TenantID,
		azureConfig.ClientID,
		azureConfig.ClientSecret,
	)
	if err != nil {
		return fmt.Errorf("error getting bearer token authorizer: %s", err)
	}

	armDeployer := arm.NewDeployer(
		azureEnvironment,
		azureConfig.SubscriptionID,
		authorizer,
	)

	aciClient := aciSDK.NewContainerGroupsClientWithBaseURI(
		azureEnvironment.ResourceManagerEndpoint,
		azureConfig.SubscriptionID,
	)
	aciClient.Authorizer = authorizer

	cosmosdbAccountsClient := cosmosSDK.NewDatabaseAccountsClientWithBaseURI(
		azureEnvironment.ResourceManagerEndpoint,
		azureConfig.SubscriptionID,
	)
	cosmosdbAccountsClient.Authorizer = authorizer

	eventHubNamespacesClient := eventHubSDK.NewNamespacesClientWithBaseURI(
		azureEnvironment.ResourceManagerEndpoint,
		azureConfig.SubscriptionID,
	)
	eventHubNamespacesClient.Authorizer = authorizer

	keyVaultsClient := keyVaultSDK.NewVaultsClientWithBaseURI(
		azureEnvironment.ResourceManagerEndpoint,
		azureConfig.SubscriptionID,
	)
	keyVaultsClient.Authorizer = authorizer

	mysqlServersClient := mysqlSDK.NewServersClientWithBaseURI(
		azureEnvironment.ResourceManagerEndpoint,
		azureConfig.SubscriptionID,
	)
	mysqlServersClient.Authorizer = authorizer

	postgresServersClient := postgresSDK.NewServersClientWithBaseURI(
		azureEnvironment.ResourceManagerEndpoint,
		azureConfig.SubscriptionID,
	)
	postgresServersClient.Authorizer = authorizer

	sqlServersClient := sqlSDK.NewServersClientWithBaseURI(
		azureEnvironment.ResourceManagerEndpoint,
		azureConfig.SubscriptionID,
	)
	sqlServersClient.Authorizer = authorizer
	sqlDatabasesClient := sqlSDK.NewDatabasesClientWithBaseURI(
		azureEnvironment.ResourceManagerEndpoint,
		azureConfig.SubscriptionID,
	)
	sqlDatabasesClient.Authorizer = authorizer

	redisGroupClient := redisSDK.NewGroupClientWithBaseURI(
		azureEnvironment.ResourceManagerEndpoint,
		azureConfig.SubscriptionID,
	)
	redisGroupClient.Authorizer = authorizer

	searchServicesClient := searchSDK.NewServicesClientWithBaseURI(
		azureEnvironment.ResourceManagerEndpoint,
		azureConfig.SubscriptionID,
	)
	searchServicesClient.Authorizer = authorizer

	serviceBusNamespacesClient := servicebusSDK.NewNamespacesClientWithBaseURI(
		azureEnvironment.ResourceManagerEndpoint,
		azureConfig.SubscriptionID,
	)
	serviceBusNamespacesClient.Authorizer = authorizer

	storageAccountsClient := storageSDK.NewAccountsClientWithBaseURI(
		azureEnvironment.ResourceManagerEndpoint,
		azureConfig.SubscriptionID,
	)
	storageAccountsClient.Authorizer = authorizer

	modules = []service.Module{
		postgresqldb.New(armDeployer, postgresServersClient),
		rediscache.New(armDeployer, redisGroupClient),
		mysqldb.New(armDeployer, mysqlServersClient),
		servicebus.New(armDeployer, serviceBusNamespacesClient),
		eventhubs.New(armDeployer, eventHubNamespacesClient),
		keyvault.New(azureConfig.TenantID, armDeployer, keyVaultsClient),
		sqldb.New(armDeployer, sqlServersClient, sqlDatabasesClient),
		cosmosdb.New(armDeployer, cosmosdbAccountsClient),
		storage.New(armDeployer, storageAccountsClient),
		search.New(armDeployer, searchServicesClient),
		aci.New(armDeployer, aciClient),
	}
	return nil
}

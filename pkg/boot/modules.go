package boot

// nolint: lll
import (
	"fmt"
	"time"

	// aciSDK "github.com/Azure/azure-sdk-for-go/services/containerinstance/mgmt/2017-08-01-preview/containerinstance"
	// cosmosSDK "github.com/Azure/azure-sdk-for-go/services/cosmos-db/mgmt/2015-04-08/documentdb"
	eventHubSDK "github.com/Azure/azure-sdk-for-go/services/eventhub/mgmt/2017-04-01/eventhub"
	// keyVaultSDK "github.com/Azure/azure-sdk-for-go/services/keyvault/mgmt/2016-10-01/keyvault"
	mysqlSDK "github.com/Azure/azure-sdk-for-go/services/mysql/mgmt/2017-04-30-preview/mysql"
	postgresSDK "github.com/Azure/azure-sdk-for-go/services/postgresql/mgmt/2017-04-30-preview/postgresql"
	// redisSDK "github.com/Azure/azure-sdk-for-go/services/redis/mgmt/2017-10-01/redis"
	resourcesSDK "github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2017-05-10/resources"
	// searchSDK "github.com/Azure/azure-sdk-for-go/services/search/mgmt/2015-08-19/search"
	// servicebusSDK "github.com/Azure/azure-sdk-for-go/services/servicebus/mgmt/2017-04-01/servicebus"
	sqlSDK "github.com/Azure/azure-sdk-for-go/services/sql/mgmt/2017-03-01-preview/sql"
	// storageSDK "github.com/Azure/azure-sdk-for-go/services/storage/mgmt/2017-10-01/storage"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/open-service-broker-azure/pkg/azure"
	"github.com/Azure/open-service-broker-azure/pkg/azure/arm"
	"github.com/Azure/open-service-broker-azure/pkg/service"
	"github.com/Azure/open-service-broker-azure/pkg/services/eventhubs"
	"github.com/Azure/open-service-broker-azure/pkg/services/mssql"
	"github.com/Azure/open-service-broker-azure/pkg/services/mysql"
	"github.com/Azure/open-service-broker-azure/pkg/services/postgresql"
	"github.com/Azure/open-service-broker-azure/pkg/version"
)

func getModules(
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

	// aciClient := aciSDK.NewContainerGroupsClientWithBaseURI(
	// 	azureConfig.Environment.ResourceManagerEndpoint,
	// 	azureSubscriptionID,
	// )
	// aciClient.Authorizer = authorizer
	// aciClient.UserAgent = getUserAgent(aciClient.Client)

	// cosmosdbAccountsClient := cosmosSDK.NewDatabaseAccountsClientWithBaseURI(
	// 	azureConfig.Environment.ResourceManagerEndpoint,
	// 	azureSubscriptionID,
	// )
	// cosmosdbAccountsClient.Authorizer = authorizer
	// cosmosdbAccountsClient.UserAgent =
	// 	getUserAgent(cosmosdbAccountsClient.Client)

	eventHubNamespacesClient := eventHubSDK.NewNamespacesClientWithBaseURI(
		azureConfig.Environment.ResourceManagerEndpoint,
		azureSubscriptionID,
	)
	eventHubNamespacesClient.Authorizer = authorizer
	eventHubNamespacesClient.UserAgent =
		getUserAgent(eventHubNamespacesClient.Client)

	// keyVaultsClient := keyVaultSDK.NewVaultsClientWithBaseURI(
	// 	azureConfig.Environment.ResourceManagerEndpoint,
	// 	azureSubscriptionID,
	// )
	// keyVaultsClient.Authorizer = authorizer
	// keyVaultsClient.UserAgent = getUserAgent(keyVaultsClient.Client)

	mysqlCheckNameAvailabilityClient :=
		mysqlSDK.NewCheckNameAvailabilityClientWithBaseURI(
			azureConfig.Environment.ResourceManagerEndpoint,
			azureSubscriptionID,
		)
	mysqlCheckNameAvailabilityClient.Authorizer = authorizer
	mysqlCheckNameAvailabilityClient.UserAgent =
		getUserAgent(mysqlCheckNameAvailabilityClient.Client)
	mysqlServersClient := mysqlSDK.NewServersClientWithBaseURI(
		azureConfig.Environment.ResourceManagerEndpoint,
		azureSubscriptionID,
	)
	mysqlServersClient.Authorizer = authorizer
	mysqlServersClient.UserAgent = getUserAgent(mysqlServersClient.Client)
	mysqlDatabasesClient := mysqlSDK.NewDatabasesClientWithBaseURI(
		azureConfig.Environment.ResourceManagerEndpoint,
		azureSubscriptionID,
	)
	mysqlDatabasesClient.Authorizer = authorizer
	mysqlDatabasesClient.UserAgent = getUserAgent(mysqlDatabasesClient.Client)

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

	sqlServersClient := sqlSDK.NewServersClientWithBaseURI(
		azureConfig.Environment.ResourceManagerEndpoint,
		azureSubscriptionID,
	)
	sqlServersClient.Authorizer = authorizer
	sqlServersClient.UserAgent = getUserAgent(sqlServersClient.Client)
	sqlDatabasesClient := sqlSDK.NewDatabasesClientWithBaseURI(
		azureConfig.Environment.ResourceManagerEndpoint,
		azureSubscriptionID,
	)
	sqlDatabasesClient.Authorizer = authorizer
	sqlDatabasesClient.UserAgent = getUserAgent(sqlDatabasesClient.Client)

	// redisClient := redisSDK.NewClientWithBaseURI(
	// 	azureConfig.Environment.ResourceManagerEndpoint,
	// 	azureSubscriptionID,
	// )
	// redisClient.Authorizer = authorizer
	// redisClient.UserAgent = getUserAgent(redisClient.Client)

	// searchServicesClient := searchSDK.NewServicesClientWithBaseURI(
	// 	azureConfig.Environment.ResourceManagerEndpoint,
	// 	azureSubscriptionID,
	// )
	// searchServicesClient.Authorizer = authorizer
	// searchServicesClient.UserAgent = getUserAgent(searchServicesClient.Client)

	// serviceBusNamespacesClient := servicebusSDK.NewNamespacesClientWithBaseURI(
	// 	azureConfig.Environment.ResourceManagerEndpoint,
	// 	azureSubscriptionID,
	// )
	// serviceBusNamespacesClient.Authorizer = authorizer
	// serviceBusNamespacesClient.UserAgent =
	// 	getUserAgent(serviceBusNamespacesClient.Client)

	// storageAccountsClient := storageSDK.NewAccountsClientWithBaseURI(
	// 	azureConfig.Environment.ResourceManagerEndpoint,
	// 	azureSubscriptionID,
	// )
	// storageAccountsClient.Authorizer = authorizer
	// storageAccountsClient.UserAgent = getUserAgent(storageAccountsClient.Client)

	modules := []service.Module{
		postgresql.New(
			armDeployer,
			postgresCheckNameAvailabilityClient,
			postgresServersClient,
			postgresDatabasesClient,
		),
		// rediscache.New(armDeployer, redisClient),
		mysql.New(
			azureConfig.Environment,
			armDeployer,
			mysqlCheckNameAvailabilityClient,
			mysqlServersClient,
			mysqlDatabasesClient,
		),
		// servicebus.New(armDeployer, serviceBusNamespacesClient),
		eventhubs.New(armDeployer, eventHubNamespacesClient),
		// keyvault.New(azureConfig.TenantID, armDeployer, keyVaultsClient),
		mssql.New(
			azureConfig.Environment,
			armDeployer,
			sqlServersClient,
			sqlDatabasesClient,
		),
		// cosmosdb.New(armDeployer, cosmosdbAccountsClient),
		// storage.New(armDeployer, storageAccountsClient),
		// search.New(armDeployer, searchServicesClient),
		// aci.New(armDeployer, aciClient),
	}

	return modules, nil
}

func getUserAgent(client autorest.Client) string {
	return fmt.Sprintf(
		"%s; open-service-broker/%s",
		client.UserAgent,
		version.GetVersion(),
	)
}

package boot

// nolint: lll
import (
	"fmt"
	"time"

	appInsightsSDK "github.com/Azure/azure-sdk-for-go/services/appinsights/mgmt/2015-05-01/insights"
	cognitiveSDK "github.com/Azure/azure-sdk-for-go/services/cognitiveservices/mgmt/2017-04-18/cognitiveservices"
	cosmosSDK "github.com/Azure/azure-sdk-for-go/services/cosmos-db/mgmt/2015-04-08/documentdb"
	eventHubSDK "github.com/Azure/azure-sdk-for-go/services/eventhub/mgmt/2017-04-01/eventhub"
	iotHubSDK "github.com/Azure/azure-sdk-for-go/services/iothub/mgmt/2017-07-01/devices"
	keyVaultSDK "github.com/Azure/azure-sdk-for-go/services/keyvault/mgmt/2016-10-01/keyvault"
	mysqlSDK "github.com/Azure/azure-sdk-for-go/services/mysql/mgmt/2017-12-01/mysql"
	postgresSDK "github.com/Azure/azure-sdk-for-go/services/postgresql/mgmt/2017-12-01/postgresql"
	sqlSDK "github.com/Azure/azure-sdk-for-go/services/preview/sql/mgmt/2017-03-01-preview/sql"
	redisSDK "github.com/Azure/azure-sdk-for-go/services/redis/mgmt/2017-10-01/redis"
	resourcesSDK "github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2017-05-10/resources"
	servicebusSDK "github.com/Azure/azure-sdk-for-go/services/servicebus/mgmt/2017-04-01/servicebus"
	storageSDK "github.com/Azure/azure-sdk-for-go/services/storage/mgmt/2017-10-01/storage"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/open-service-broker-azure/pkg/azure"
	"github.com/Azure/open-service-broker-azure/pkg/azure/arm"
	"github.com/Azure/open-service-broker-azure/pkg/service"
	"github.com/Azure/open-service-broker-azure/pkg/services/appinsights"
	"github.com/Azure/open-service-broker-azure/pkg/services/cosmosdb"
	"github.com/Azure/open-service-broker-azure/pkg/services/eventhubs"
	"github.com/Azure/open-service-broker-azure/pkg/services/iothub"
	"github.com/Azure/open-service-broker-azure/pkg/services/keyvault"
	"github.com/Azure/open-service-broker-azure/pkg/services/mssql"
	"github.com/Azure/open-service-broker-azure/pkg/services/mssqldr"
	"github.com/Azure/open-service-broker-azure/pkg/services/mysql"
	"github.com/Azure/open-service-broker-azure/pkg/services/postgresql"
	"github.com/Azure/open-service-broker-azure/pkg/services/rediscache"
	"github.com/Azure/open-service-broker-azure/pkg/services/servicebus"
	"github.com/Azure/open-service-broker-azure/pkg/services/storage"
	"github.com/Azure/open-service-broker-azure/pkg/services/textanalytics"
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
	resourceDeploymentsClient.PollingDuration = time.Minute * 45
	armDeployer := arm.NewDeployer(
		resourceGroupsClient,
		resourceDeploymentsClient,
	)

	appInsightsClient := appInsightsSDK.NewComponentsClientWithBaseURI(
		azureConfig.Environment.ResourceManagerEndpoint,
		azureSubscriptionID,
	)
	appInsightsClient.Authorizer = authorizer
	appInsightsClient.UserAgent = getUserAgent(appInsightsClient.Client)
	appInsightsAPIKeyClient := appInsightsSDK.NewAPIKeysClientWithBaseURI(
		azureConfig.Environment.ResourceManagerEndpoint,
		azureSubscriptionID,
	)
	appInsightsAPIKeyClient.Authorizer = authorizer
	appInsightsAPIKeyClient.UserAgent = getUserAgent(
		appInsightsAPIKeyClient.Client,
	)

	cognitiveClient := cognitiveSDK.NewAccountsClientWithBaseURI(
		azureConfig.Environment.ResourceManagerEndpoint,
		azureSubscriptionID,
	)
	cognitiveClient.Authorizer = authorizer
	cognitiveClient.UserAgent = getUserAgent(cognitiveClient.Client)

	cosmosdbAccountsClient := cosmosSDK.NewDatabaseAccountsClientWithBaseURI(
		azureConfig.Environment.ResourceManagerEndpoint,
		azureSubscriptionID,
	)
	cosmosdbAccountsClient.Authorizer = authorizer
	cosmosdbAccountsClient.UserAgent =
		getUserAgent(cosmosdbAccountsClient.Client)
	// When there are multiple read regions, default polling duration
	// (15 minutes) is not enough for deleting all of the regions.
	cosmosdbAccountsClient.PollingDuration = time.Minute * 30

	eventHubNamespacesClient := eventHubSDK.NewNamespacesClientWithBaseURI(
		azureConfig.Environment.ResourceManagerEndpoint,
		azureSubscriptionID,
	)
	eventHubNamespacesClient.Authorizer = authorizer
	eventHubNamespacesClient.UserAgent =
		getUserAgent(eventHubNamespacesClient.Client)

	iotHubClient := iotHubSDK.NewIotHubResourceClientWithBaseURI(
		azureConfig.Environment.ResourceManagerEndpoint,
		azureSubscriptionID,
	)
	iotHubClient.Authorizer = authorizer
	iotHubClient.UserAgent =
		getUserAgent(iotHubClient.Client)

	keyVaultsClient := keyVaultSDK.NewVaultsClientWithBaseURI(
		azureConfig.Environment.ResourceManagerEndpoint,
		azureSubscriptionID,
	)
	keyVaultsClient.Authorizer = authorizer
	keyVaultsClient.UserAgent = getUserAgent(keyVaultsClient.Client)

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
	sqlServerConnectionPoliciesClient :=
		sqlSDK.NewServerConnectionPoliciesClientWithBaseURI(
			azureConfig.Environment.ResourceManagerEndpoint,
			azureSubscriptionID,
		)
	sqlServerConnectionPoliciesClient.Authorizer = authorizer
	sqlServerConnectionPoliciesClient.UserAgent = getUserAgent(
		sqlServerConnectionPoliciesClient.Client,
	)

	sqlFailoverGroupsClient := sqlSDK.NewFailoverGroupsClientWithBaseURI(
		azureConfig.Environment.ResourceManagerEndpoint,
		azureSubscriptionID,
	)
	sqlFailoverGroupsClient.Authorizer = authorizer
	sqlFailoverGroupsClient.UserAgent =
		getUserAgent(sqlFailoverGroupsClient.Client)

	redisClient := redisSDK.NewClientWithBaseURI(
		azureConfig.Environment.ResourceManagerEndpoint,
		azureSubscriptionID,
	)
	redisClient.Authorizer = authorizer
	redisClient.UserAgent = getUserAgent(redisClient.Client)

	serviceBusNamespacesClient := servicebusSDK.NewNamespacesClientWithBaseURI(
		azureConfig.Environment.ResourceManagerEndpoint,
		azureSubscriptionID,
	)
	serviceBusNamespacesClient.Authorizer = authorizer
	serviceBusNamespacesClient.UserAgent =
		getUserAgent(serviceBusNamespacesClient.Client)
	serviceBusQueuesClient := servicebusSDK.NewQueuesClientWithBaseURI(
		azureConfig.Environment.ResourceManagerEndpoint,
		azureSubscriptionID,
	)
	serviceBusQueuesClient.Authorizer = authorizer
	serviceBusQueuesClient.UserAgent =
		getUserAgent(serviceBusQueuesClient.Client)
	serviceBusTopicsClient := servicebusSDK.NewTopicsClientWithBaseURI(
		azureConfig.Environment.ResourceManagerEndpoint,
		azureSubscriptionID,
	)
	serviceBusTopicsClient.Authorizer = authorizer
	serviceBusTopicsClient.UserAgent =
		getUserAgent(serviceBusTopicsClient.Client)
	serviceBusSubscriptionsClient := servicebusSDK.NewSubscriptionsClientWithBaseURI( // nolint: lll
		azureConfig.Environment.ResourceManagerEndpoint,
		azureSubscriptionID,
	)
	serviceBusSubscriptionsClient.Authorizer = authorizer
	serviceBusSubscriptionsClient.UserAgent =
		getUserAgent(serviceBusQueuesClient.Client)

	storageAccountsClient := storageSDK.NewAccountsClientWithBaseURI(
		azureConfig.Environment.ResourceManagerEndpoint,
		azureSubscriptionID,
	)
	storageAccountsClient.Authorizer = authorizer
	storageAccountsClient.UserAgent = getUserAgent(storageAccountsClient.Client)

	modules := []service.Module{
		postgresql.New(
			armDeployer,
			postgresCheckNameAvailabilityClient,
			postgresServersClient,
			postgresDatabasesClient,
		),
		rediscache.New(armDeployer, redisClient),
		mysql.New(
			azureConfig.Environment,
			armDeployer,
			mysqlCheckNameAvailabilityClient,
			mysqlServersClient,
			mysqlDatabasesClient,
		),
		servicebus.New(
			armDeployer,
			serviceBusNamespacesClient,
			serviceBusQueuesClient,
			serviceBusTopicsClient,
			serviceBusSubscriptionsClient,
		),
		eventhubs.New(armDeployer, eventHubNamespacesClient),
		keyvault.New(azureConfig.TenantID, armDeployer, keyVaultsClient),
		mssql.New(
			azureConfig.Environment,
			armDeployer,
			sqlServersClient,
			sqlDatabasesClient,
			sqlServerConnectionPoliciesClient,
		),
		mssqldr.New(
			azureConfig.Environment,
			armDeployer,
			sqlServersClient,
			sqlDatabasesClient,
			sqlFailoverGroupsClient,
		),
		cosmosdb.New(armDeployer, cosmosdbAccountsClient),
		storage.New(armDeployer, storageAccountsClient),
		textanalytics.New(armDeployer, cognitiveClient),
		iothub.New(armDeployer, iotHubClient),
		appinsights.New(armDeployer, appInsightsClient, appInsightsAPIKeyClient),
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

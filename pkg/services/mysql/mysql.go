package mysql

import (
	mysqlSDK "github.com/Azure/azure-sdk-for-go/services/mysql/mgmt/2017-04-30-preview/mysql" // nolint: lll
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/open-service-broker-azure/pkg/azure/arm"
	"github.com/Azure/open-service-broker-azure/pkg/service"
)

type module struct {
	allInOneServiceManager *allInOneManager
	dbmsManager            *dbmsManager
	databaseManager        *databaseManager
}

type allInOneManager struct {
	sqlDatabaseDNSSuffix        string
	armDeployer                 arm.Deployer
	checkNameAvailabilityClient mysqlSDK.CheckNameAvailabilityClient
	serversClient               mysqlSDK.ServersClient
}

type dbmsManager struct {
	sqlDatabaseDNSSuffix        string
	armDeployer                 arm.Deployer
	checkNameAvailabilityClient mysqlSDK.CheckNameAvailabilityClient
	serversClient               mysqlSDK.ServersClient
}

type databaseManager struct {
	sqlDatabaseDNSSuffix string
	armDeployer          arm.Deployer
	databasesClient      mysqlSDK.DatabasesClient
}

// New returns a new instance of a type that fulfills the service.Module
// interface and is capable of provisioning MySQL DBMS and databases
// using "Azure Database for MySQL"
func New(
	azureEnvironment azure.Environment,
	armDeployer arm.Deployer,
	checkNameAvailabilityClient mysqlSDK.CheckNameAvailabilityClient,
	serversClient mysqlSDK.ServersClient,
	databaseClient mysqlSDK.DatabasesClient,
) service.Module {
	return &module{
		allInOneServiceManager: &allInOneManager{
			sqlDatabaseDNSSuffix:        azureEnvironment.SQLDatabaseDNSSuffix,
			armDeployer:                 armDeployer,
			checkNameAvailabilityClient: checkNameAvailabilityClient,
			serversClient:               serversClient,
		},
		dbmsManager: &dbmsManager{
			sqlDatabaseDNSSuffix:        azureEnvironment.SQLDatabaseDNSSuffix,
			armDeployer:                 armDeployer,
			checkNameAvailabilityClient: checkNameAvailabilityClient,
			serversClient:               serversClient,
		},
		databaseManager: &databaseManager{
			sqlDatabaseDNSSuffix: azureEnvironment.SQLDatabaseDNSSuffix,
			armDeployer:          armDeployer,
			databasesClient:      databaseClient,
		},
	}
}

func (m *module) GetName() string {
	return "mysql"
}

func (m *module) GetStability() service.Stability {
	return service.StabilityExperimental
}

package mysqldb

import (
	mysqlSDK "github.com/Azure/azure-sdk-for-go/services/mysql/mgmt/2017-04-30-preview/mysql" // nolint: lll
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/open-service-broker-azure/pkg/azure/arm"
	"github.com/Azure/open-service-broker-azure/pkg/service"
)

type module struct {
	allInOneServiceManager *allInOneManager
	dbmsOnlyManager        *dbmsOnlyManager
	dbOnlyServiceManager   *dbOnlyManager
}

type allInOneManager struct {
	sqlDatabaseDNSSuffix string
	armDeployer          arm.Deployer
	serversClient        mysqlSDK.ServersClient
}

type dbmsOnlyManager struct {
	sqlDatabaseDNSSuffix string
	armDeployer          arm.Deployer
	serversClient        mysqlSDK.ServersClient
}

type dbOnlyManager struct {
	sqlDatabaseDNSSuffix string
	armDeployer          arm.Deployer
	databasesClient      mysqlSDK.DatabasesClient
}

// New returns a new instance of a type that fulfills the service.Module
// interface and is capable of provisioning MySQL servers and databases
// using "Azure Database for MySQL"
func New(
	azureEnvironment azure.Environment,
	armDeployer arm.Deployer,
	serversClient mysqlSDK.ServersClient,
	databaseClient mysqlSDK.DatabasesClient,
) service.Module {
	return &module{
		allInOneServiceManager: &allInOneManager{
			sqlDatabaseDNSSuffix: azureEnvironment.SQLDatabaseDNSSuffix,
			armDeployer:          armDeployer,
			serversClient:        serversClient,
		},
		dbmsOnlyManager: &dbmsOnlyManager{
			sqlDatabaseDNSSuffix: azureEnvironment.SQLDatabaseDNSSuffix,
			armDeployer:          armDeployer,
			serversClient:        serversClient,
		},
		dbOnlyServiceManager: &dbOnlyManager{
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

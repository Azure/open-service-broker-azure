package mssql

import (
	sqlSDK "github.com/Azure/azure-sdk-for-go/services/sql/mgmt/2017-03-01-preview/sql" // nolint: lll
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/open-service-broker-azure/pkg/azure/arm"
	"github.com/Azure/open-service-broker-azure/pkg/service"
)

type module struct {
	allInOneServiceManager *allInOneManager
	dbmsManager            *dbmsManager
	databaseManager        *databaseManager
	dbmsRegisteredManager  *dbmsRegisteredManager
	databaseFeManager      *databaseFeManager
}

type allInOneManager struct {
	sqlDatabaseDNSSuffix string
	armDeployer          arm.Deployer
	serversClient        sqlSDK.ServersClient
	databasesClient      sqlSDK.DatabasesClient
}

type dbmsManager struct {
	sqlDatabaseDNSSuffix string
	armDeployer          arm.Deployer
	serversClient        sqlSDK.ServersClient
}

type databaseManager struct {
	armDeployer     arm.Deployer
	databasesClient sqlSDK.DatabasesClient
}

type dbmsRegisteredManager struct {
	dbmsManager
}

type databaseFeManager struct {
	databaseManager
}

// New returns a new instance of a type that fulfills the service.Module
// interface and is capable of provisioning MS SQL servers and databases
// using "Azure SQL Database"
func New(
	azureEnvironment azure.Environment,
	armDeployer arm.Deployer,
	serversClient sqlSDK.ServersClient,
	databasesClient sqlSDK.DatabasesClient,
) service.Module {
	dbmsMan := dbmsManager{
		sqlDatabaseDNSSuffix: azureEnvironment.SQLDatabaseDNSSuffix,
		armDeployer:          armDeployer,
		serversClient:        serversClient,
	}
	databaseMan := databaseManager{
		armDeployer:     armDeployer,
		databasesClient: databasesClient,
	}
	return &module{
		allInOneServiceManager: &allInOneManager{
			sqlDatabaseDNSSuffix: azureEnvironment.SQLDatabaseDNSSuffix,
			armDeployer:          armDeployer,
			serversClient:        serversClient,
			databasesClient:      databasesClient,
		},
		dbmsManager:     &dbmsMan,
		databaseManager: &databaseMan,
		dbmsRegisteredManager: &dbmsRegisteredManager{
			dbmsMan,
		},
		databaseFeManager: &databaseFeManager{
			databaseMan,
		},
	}
}

func (m *module) GetName() string {
	return "mssql"
}

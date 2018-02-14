package sqldb

import (
	sqlSDK "github.com/Azure/azure-sdk-for-go/arm/sql"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/open-service-broker-azure/pkg/azure/arm"
	"github.com/Azure/open-service-broker-azure/pkg/service"
)

type module struct {
	allInOneServiceManager *allInOneManager
	vmOnlyServiceManager   *vmOnlyManager
	dbOnlyServiceManager   *dbOnlyManager
}

type allInOneManager struct {
	armDeployer     arm.Deployer
	serversClient   sqlSDK.ServersClient
	databasesClient sqlSDK.DatabasesClient
}

type vmOnlyManager struct {
	armDeployer   arm.Deployer
	serversClient sqlSDK.ServersClient
}

type dbOnlyManager struct {
	sqlDatabaseDNSSuffix string
	armDeployer          arm.Deployer
	databasesClient      sqlSDK.DatabasesClient
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
	return &module{
		allInOneServiceManager: &allInOneManager{
			armDeployer:     armDeployer,
			serversClient:   serversClient,
			databasesClient: databasesClient,
		},
		vmOnlyServiceManager: &vmOnlyManager{
			armDeployer:   armDeployer,
			serversClient: serversClient,
		},
		dbOnlyServiceManager: &dbOnlyManager{
			sqlDatabaseDNSSuffix: azureEnvironment.SQLDatabaseDNSSuffix,
			armDeployer:          armDeployer,
			databasesClient:      databasesClient,
		},
	}
}

func (m *module) GetName() string {
	return "mssql"
}

func (m *module) GetStability() service.Stability {
	return service.StabilityExperimental
}

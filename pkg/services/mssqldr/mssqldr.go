package mssqldr

import (
	sqlSDK "github.com/Azure/azure-sdk-for-go/services/sql/mgmt/2017-03-01-preview/sql" // nolint: lll
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/open-service-broker-azure/pkg/azure/arm"
	"github.com/Azure/open-service-broker-azure/pkg/service"
)

type module struct {
	dbmsPairRegisteredManager     *dbmsPairRegisteredManager
	databasePairManager           *databasePairManager
	databasePairRegisteredManager *databasePairRegisteredManager
}

type dbmsPairRegisteredManager struct {
	sqlDatabaseDNSSuffix string
	armDeployer          arm.Deployer
	serversClient        sqlSDK.ServersClient
}

type commonDatabasePairManager struct {
	armDeployer          arm.Deployer
	databasesClient      sqlSDK.DatabasesClient
	failoverGroupsClient sqlSDK.FailoverGroupsClient
}

type databasePairManager struct {
	commonDatabasePairManager
}

type databasePairRegisteredManager struct {
	commonDatabasePairManager
}

// New returns a new instance of a type that fulfills the service.Module
// interface and is capable of provisioning MS SQL servers and databases
// using "Azure SQL Database"
func New(
	azureEnvironment azure.Environment,
	armDeployer arm.Deployer,
	serversClient sqlSDK.ServersClient,
	databasesClient sqlSDK.DatabasesClient,
	failoverGroupsClient sqlSDK.FailoverGroupsClient,
) service.Module {
	commonDatabasePairMgr := commonDatabasePairManager{
		armDeployer:          armDeployer,
		databasesClient:      databasesClient,
		failoverGroupsClient: failoverGroupsClient,
	}
	return &module{
		dbmsPairRegisteredManager: &dbmsPairRegisteredManager{
			sqlDatabaseDNSSuffix: azureEnvironment.SQLDatabaseDNSSuffix,
			armDeployer:          armDeployer,
			serversClient:        serversClient,
		},
		databasePairManager: &databasePairManager{
			commonDatabasePairMgr,
		},
		databasePairRegisteredManager: &databasePairRegisteredManager{
			commonDatabasePairMgr,
		},
	}
}

func (m *module) GetName() string {
	return "mssqldr"
}

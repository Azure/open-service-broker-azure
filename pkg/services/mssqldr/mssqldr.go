package mssqldr

import (
	sqlSDK "github.com/Azure/azure-sdk-for-go/services/sql/mgmt/2017-03-01-preview/sql" // nolint: lll
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/open-service-broker-azure/pkg/azure/arm"
	"github.com/Azure/open-service-broker-azure/pkg/service"
)

type module struct {
	dbmsPairRegisteredManager             *dbmsPairRegisteredManager
	databasePairManager                   *databasePairManager
	databasePairRegisteredManager         *databasePairRegisteredManager
	databasePairManagerForExistingPrimary *databasePairManagerForExistingPrimary
	databasePairManagerForExistingPair    *databasePairManagerForExistingPair
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

type databasePairManagerForExistingPrimary struct {
	commonDatabasePairManager
}

type databasePairManagerForExistingPair struct {
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
		databasePairManagerForExistingPrimary: &databasePairManagerForExistingPrimary{ // nolint: lll
			commonDatabasePairMgr,
		},
		databasePairManagerForExistingPair: &databasePairManagerForExistingPair{ // nolint: lll
			commonDatabasePairMgr,
		},
	}
}

func (m *module) GetName() string {
	return "mssqldr"
}

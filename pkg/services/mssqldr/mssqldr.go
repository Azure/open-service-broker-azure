package mssqldr

import (
	sqlSDK "github.com/Azure/azure-sdk-for-go/services/sql/mgmt/2017-03-01-preview/sql" // nolint: lll
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/open-service-broker-azure/pkg/azure/arm"
	"github.com/Azure/open-service-broker-azure/pkg/service"
)

type module struct {
	dbmsPairRegisteredManager *dbmsPairRegisteredManager
}

// The basic parent service in this module. A pair of servers is required to
// foster the foundation of a Failover Group. The reason why not new server
// pair, is based on users' feedback. We can still add a service for new server
// pair if needed.
type dbmsPairRegisteredManager struct {
	sqlDatabaseDNSSuffix string
	armDeployer          arm.Deployer
	serversClient        sqlSDK.ServersClient
}

// New returns a new instance of a type that fulfills the service.Module
// interface and is capable of provisioning MS SQL servers and databases
// using "Azure SQL Database"
func New(
	azureEnvironment azure.Environment,
	armDeployer arm.Deployer,
	serversClient sqlSDK.ServersClient,
	databasesClient sqlSDK.DatabasesClient, // nolint: unparam
	failoverGroupsClient sqlSDK.FailoverGroupsClient, // nolint: unparam
	// These used clients are for future PRs usage
) service.Module {
	return &module{
		dbmsPairRegisteredManager: &dbmsPairRegisteredManager{
			sqlDatabaseDNSSuffix: azureEnvironment.SQLDatabaseDNSSuffix,
			armDeployer:          armDeployer,
			serversClient:        serversClient,
		},
	}
}

func (m *module) GetName() string {
	return "mssqldr"
}

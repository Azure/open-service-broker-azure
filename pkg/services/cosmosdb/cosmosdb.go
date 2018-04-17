package cosmosdb

import (
	cosmosSDK "github.com/Azure/azure-sdk-for-go/services/cosmos-db/mgmt/2015-04-08/documentdb" // nolint: lll
	"github.com/Azure/open-service-broker-azure/pkg/azure/arm"
	"github.com/Azure/open-service-broker-azure/pkg/service"
)

type module struct {
	sqlAccountManager   *sqlAccountManager
	sqlAllInOneManager  *sqlAllInOneManager
	mongoAccountManager *mongoAccountManager
	graphAccountManager *graphAccountManager
	tableAccountManager *tableAccountManager
}

type cosmosAccountManager struct {
	armDeployer            arm.Deployer
	databaseAccountsClient cosmosSDK.DatabaseAccountsClient
}

type sqlAccountManager struct {
	cosmosAccountManager
}

type sqlAllInOneManager struct {
	sqlAccountManager
}

type mongoAccountManager struct {
	cosmosAccountManager
}

type tableAccountManager struct {
	cosmosAccountManager
}

type graphAccountManager struct {
	cosmosAccountManager
}

// New returns a new instance of a type that fulfills the service.Module
// interface and is capable of provisioning CosmosDB database accounts and
// databases using "Azure Database for CosmosDB"
func New(
	armDeployer arm.Deployer,
	databaseAccountsClient cosmosSDK.DatabaseAccountsClient,
) service.Module {
	cosmos := cosmosAccountManager{
		armDeployer:            armDeployer,
		databaseAccountsClient: databaseAccountsClient,
	}
	return &module{
		mongoAccountManager: &mongoAccountManager{cosmos},
		sqlAccountManager:   &sqlAccountManager{cosmos},
		sqlAllInOneManager: &sqlAllInOneManager{
			sqlAccountManager: sqlAccountManager{cosmos},
		},
		graphAccountManager: &graphAccountManager{cosmos},
		tableAccountManager: &tableAccountManager{cosmos},
	}
}

func (m *module) GetName() string {
	return "cosmosdb"
}

func (m *module) GetStability() service.Stability {
	return service.StabilityExperimental
}

package cosmosdb

import (
	cosmosSDK "github.com/Azure/azure-sdk-for-go/services/cosmos-db/mgmt/2015-04-08/documentdb" // nolint: lll
	"github.com/Azure/open-service-broker-azure/pkg/azure/arm"
	"github.com/Azure/open-service-broker-azure/pkg/service"
)

type module struct {
	cosmosAccountManager *cosmosAccountManager
	mongoAccountManager  *mongoAccountManager
}

type cosmosAccountManager struct {
	armDeployer            arm.Deployer
	databaseAccountsClient cosmosSDK.DatabaseAccountsClient
}

type mongoAccountManager struct {
	armDeployer            arm.Deployer
	databaseAccountsClient cosmosSDK.DatabaseAccountsClient
}

// New returns a new instance of a type that fulfills the service.Module
// interface and is capable of provisioning CosmosDB database accounts and
// databases using "Azure Database for CosmosDB"
func New(
	armDeployer arm.Deployer,
	databaseAccountsClient cosmosSDK.DatabaseAccountsClient,
) service.Module {
	return &module{
		mongoAccountManager: &mongoAccountManager{
			armDeployer:            armDeployer,
			databaseAccountsClient: databaseAccountsClient,
		},
		cosmosAccountManager: &cosmosAccountManager{
			armDeployer:            armDeployer,
			databaseAccountsClient: databaseAccountsClient,
		},
	}
}

func (m *module) GetName() string {
	return "cosmosdb"
}

func (m *module) GetStability() service.Stability {
	return service.StabilityExperimental
}

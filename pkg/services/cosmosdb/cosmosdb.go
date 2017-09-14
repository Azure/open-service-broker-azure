package cosmosdb

import (
	"github.com/Azure/azure-service-broker/pkg/azure/arm"
	"github.com/Azure/azure-service-broker/pkg/azure/cosmosdb"
	"github.com/Azure/azure-service-broker/pkg/service"
)

type module struct {
	armDeployer     arm.Deployer
	cosmosdbManager cosmosdb.Manager
}

// New returns a new instance of a type that fulfills the service.Module
// interface and is capable of provisioning CosmosDB servers and databases
// using "Azure Database for CosmosDB"
func New(
	armDeployer arm.Deployer,
	cosmosdbManager cosmosdb.Manager,
) service.Module {
	return &module{
		armDeployer:     armDeployer,
		cosmosdbManager: cosmosdbManager,
	}
}

func (m *module) GetName() string {
	return "cosmosdb"
}

func (m *module) GetStability() service.Stability {
	return service.StabilityAlpha
}

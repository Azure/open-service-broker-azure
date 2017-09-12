package postgresql

import (
	"github.com/Azure/azure-service-broker/pkg/azure/arm"
	"github.com/Azure/azure-service-broker/pkg/azure/postgresql"
	"github.com/Azure/azure-service-broker/pkg/service"
)

type module struct {
	armDeployer       arm.Deployer
	postgresqlManager postgresql.Manager
}

// New returns a new instance of a type that fulfills the service.Module
// interface and is capable of provisioning PostgreSQL servers and databases
// using "Azure Database for PostgreSQL"
func New(
	armDeployer arm.Deployer,
	postgresqlManager postgresql.Manager,
) service.Module {
	return &module{
		armDeployer:       armDeployer,
		postgresqlManager: postgresqlManager,
	}
}

func (m *module) GetName() string {
	return "postgresql"
}

func (m *module) GetStability() service.Stability {
	return service.StabilityAlpha
}

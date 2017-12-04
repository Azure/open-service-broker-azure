package postgresqldb

import (
	"github.com/Azure/open-service-broker-azure/pkg/azure/arm"
	"github.com/Azure/open-service-broker-azure/pkg/azure/postgresql"
	"github.com/Azure/open-service-broker-azure/pkg/service"
)

type module struct {
	serviceManager *serviceManager
}

type serviceManager struct {
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
		serviceManager: &serviceManager{
			armDeployer:       armDeployer,
			postgresqlManager: postgresqlManager,
		},
	}
}

func (m *module) GetName() string {
	return "postgresql"
}

func (m *module) GetStability() service.Stability {
	return service.StabilityExperimental
}

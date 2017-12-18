package sqldb

import (
	"github.com/Azure/open-service-broker-azure/pkg/azure/arm"
	"github.com/Azure/open-service-broker-azure/pkg/azure/mssql"
	"github.com/Azure/open-service-broker-azure/pkg/service"
)

type module struct {
	serviceManager *serviceManager
}

type serviceManager struct {
	armDeployer  arm.Deployer
	mssqlManager mssql.Manager
}

// New returns a new instance of a type that fulfills the service.Module
// interface and is capable of provisioning MS SQL servers and databases
// using "Azure SQL Database"
func New(
	armDeployer arm.Deployer,
	mssqlManager mssql.Manager,
) service.Module {
	return &module{
		serviceManager: &serviceManager{
			armDeployer:  armDeployer,
			mssqlManager: mssqlManager,
		},
	}
}

func (m *module) GetName() string {
	return "mssql"
}

func (m *module) GetStability() service.Stability {
	return service.StabilityExperimental
}

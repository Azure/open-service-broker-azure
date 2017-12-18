package sqldb

import (
	"github.com/Azure/open-service-broker-azure/pkg/azure/arm"
	"github.com/Azure/open-service-broker-azure/pkg/azure/mssql"
	"github.com/Azure/open-service-broker-azure/pkg/service"
)

type module struct {
	allInOneServiceManager   *allInOneServiceManager
	serverOnlyServiceManager *serverOnlyServiceManager
}

type serviceManager struct {
	armDeployer  arm.Deployer
	mssqlManager mssql.Manager
}

type allInOneServiceManager struct {
	serviceManager
}

type serverOnlyServiceManager struct {
	serviceManager
}

// New returns a new instance of a type that fulfills the service.Module
// interface and is capable of provisioning MS SQL servers and databases
// using "Azure SQL Database"
func New(
	armDeployer arm.Deployer,
	mssqlManager mssql.Manager,
) service.Module {
	sm := serviceManager{
		armDeployer:  armDeployer,
		mssqlManager: mssqlManager,
	}

	return &module{
		allInOneServiceManager: &allInOneServiceManager{
			serviceManager: sm,
		},
		serverOnlyServiceManager: &serverOnlyServiceManager{
			serviceManager: sm,
		},
	}
}

func (m *module) GetName() string {
	return "mssql"
}

func (m *module) GetStability() service.Stability {
	return service.StabilityExperimental
}

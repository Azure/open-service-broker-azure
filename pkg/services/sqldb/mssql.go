package sqldb

import (
	"github.com/Azure/open-service-broker-azure/pkg/azure/arm"
	"github.com/Azure/open-service-broker-azure/pkg/azure/mssql"
	"github.com/Azure/open-service-broker-azure/pkg/service"
)

type module struct {
	allInOneServiceManager *allInOneManager
	vmOnlyServiceManager   *vmOnlyManager
	dbOnlyServiceManager   *dbOnlyManager
}

type allInOneManager struct {
	armDeployer  arm.Deployer
	mssqlManager mssql.Manager
}

type vmOnlyManager struct {
	armDeployer  arm.Deployer
	mssqlManager mssql.Manager
}

type dbOnlyManager struct {
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
		allInOneServiceManager: &allInOneManager{
			armDeployer:  armDeployer,
			mssqlManager: mssqlManager,
		},
		vmOnlyServiceManager: &vmOnlyManager{
			armDeployer:  armDeployer,
			mssqlManager: mssqlManager,
		},
		dbOnlyServiceManager: &dbOnlyManager{
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

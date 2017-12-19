package sqldb

import (
	"github.com/Azure/open-service-broker-azure/pkg/azure/arm"
	"github.com/Azure/open-service-broker-azure/pkg/azure/mssql"
	"github.com/Azure/open-service-broker-azure/pkg/service"
)

//module contains three service managers
//to refelct the three services available
//from sqldb.
type module struct {
	allInOneServiceManager *allServiceManager
	vmOnlyServiceManager   *vmServiceManager
	dbOnlyServiceManager   *dbServiceManager
}

//the default service manager for the all-in-one case
type allServiceManager struct {
	armDeployer  arm.Deployer
	mssqlManager mssql.Manager
}

//the service manager for the server only case
type vmServiceManager struct {
	armDeployer  arm.Deployer
	mssqlManager mssql.Manager
}

//the service manager for the db only case
type dbServiceManager struct {
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
		allInOneServiceManager: &allServiceManager{
			armDeployer:  armDeployer,
			mssqlManager: mssqlManager,
		},
		vmOnlyServiceManager: &vmServiceManager{
			armDeployer:  armDeployer,
			mssqlManager: mssqlManager,
		},
		dbOnlyServiceManager: &dbServiceManager{
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

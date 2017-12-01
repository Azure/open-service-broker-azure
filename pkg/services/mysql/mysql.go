package mysql

import (
	"github.com/Azure/azure-service-broker/pkg/azure/arm"
	"github.com/Azure/azure-service-broker/pkg/azure/mysql"
	"github.com/Azure/azure-service-broker/pkg/service"
)

type module struct {
	serviceManager *serviceManager
}

type serviceManager struct {
	armDeployer  arm.Deployer
	mysqlManager mysql.Manager
}

// New returns a new instance of a type that fulfills the service.Module
// interface and is capable of provisioning MySQL servers and databases
// using "Azure Database for MySQL"
func New(
	armDeployer arm.Deployer,
	mysqlManager mysql.Manager,
) service.Module {
	return &module{
		serviceManager: &serviceManager{
			armDeployer:  armDeployer,
			mysqlManager: mysqlManager,
		},
	}
}

func (m *module) GetName() string {
	return "mysql"
}

func (m *module) GetStability() service.Stability {
	return service.StabilityAlpha
}

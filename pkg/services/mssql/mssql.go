package mssql

import (
	"github.com/Azure/azure-service-broker/pkg/azure/arm"
	"github.com/Azure/azure-service-broker/pkg/azure/mssql"
	"github.com/Azure/azure-service-broker/pkg/service"
)

type module struct {
	armDeployer  arm.Deployer
	mssqlManager mssql.Manager
	mssqlConfig  Config
}

// New returns a new instance of a type that fulfills the service.Module
// interface and is capable of provisioning MS SQL servers and databases
// using "Azure SQL Database"
func New(
	armDeployer arm.Deployer,
	mssqlManager mssql.Manager,
	mssqlConfig Config,
) service.Module {
	return &module{
		armDeployer:  armDeployer,
		mssqlManager: mssqlManager,
		mssqlConfig:  mssqlConfig,
	}
}

func (m *module) GetName() string {
	return "mssql"
}

func (m *module) GetStability() service.Stability {
	return service.StabilityAlpha
}

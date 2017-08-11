package main

import (
	"fmt"

	"github.com/Azure/azure-service-broker/pkg/azure/arm"
	pg "github.com/Azure/azure-service-broker/pkg/azure/postgresql"
	"github.com/Azure/azure-service-broker/pkg/service"
	"github.com/Azure/azure-service-broker/pkg/services/echo"
	"github.com/Azure/azure-service-broker/pkg/services/postgresql"
)

var modules []service.Module

func initModules() error {
	armDeployer, err := arm.NewDeployer()
	if err != nil {
		return fmt.Errorf("error initializing ARM template deployer: %s", err)
	}
	postgreSQLManager, err := pg.NewManager()
	if err != nil {
		return fmt.Errorf("error initializing postgresql manager: %s", err)
	}
	modules = []service.Module{
		echo.New(),
		postgresql.New(armDeployer, postgreSQLManager),
	}
	return nil
}

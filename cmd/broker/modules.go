package main

import (
	"fmt"

	"github.com/Azure/azure-service-broker/pkg/azure/arm"
	mg "github.com/Azure/azure-service-broker/pkg/azure/mysql"
	pg "github.com/Azure/azure-service-broker/pkg/azure/postgresql"
	"github.com/Azure/azure-service-broker/pkg/service"
	"github.com/Azure/azure-service-broker/pkg/services/echo"
	"github.com/Azure/azure-service-broker/pkg/services/mysql"
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
	mySQLManager, err := mg.NewManager()
	if err != nil {
		return fmt.Errorf("error initializing mysql manager: %s", err)
	}

	modules = []service.Module{
		echo.New(),
		postgresql.New(armDeployer, postgreSQLManager),
		mysql.New(armDeployer, mySQLManager),
	}
	return nil
}

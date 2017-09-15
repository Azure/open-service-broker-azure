package main

import (
	"fmt"

	"github.com/Azure/azure-service-broker/pkg/azure/arm"
	cd "github.com/Azure/azure-service-broker/pkg/azure/cosmosdb"
	ss "github.com/Azure/azure-service-broker/pkg/azure/mssql"
	mg "github.com/Azure/azure-service-broker/pkg/azure/mysql"
	pg "github.com/Azure/azure-service-broker/pkg/azure/postgresql"
	rc "github.com/Azure/azure-service-broker/pkg/azure/rediscache"

	"github.com/Azure/azure-service-broker/pkg/service"
	"github.com/Azure/azure-service-broker/pkg/services/cosmosdb"
	"github.com/Azure/azure-service-broker/pkg/services/mssql"
	"github.com/Azure/azure-service-broker/pkg/services/mysql"
	"github.com/Azure/azure-service-broker/pkg/services/postgresql"
	"github.com/Azure/azure-service-broker/pkg/services/rediscache"
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
	redisManager, err := rc.NewManager()
	if err != nil {
		return fmt.Errorf("error initializing redis manager: %s", err)
	}
	msSQLManager, err := ss.NewManager()
	if err != nil {
		return fmt.Errorf("error initializing mssql manager: %s", err)
	}
	msSQLConfig, err := mssql.GetConfig()
	if err != nil {
		return fmt.Errorf("error parsing mssql configuration: %s", err)
	}
	cosmosDBManager, err := cd.NewManager()
	if err != nil {
		return fmt.Errorf("error initializing cosmosdb manager: %s", err)
	}
	modules = []service.Module{
		postgresql.New(armDeployer, postgreSQLManager),
		rediscache.New(armDeployer, redisManager),
		mysql.New(armDeployer, mySQLManager),
		mssql.New(armDeployer, msSQLManager, msSQLConfig),
		cosmosdb.New(armDeployer, cosmosDBManager),
	}
	return nil
}

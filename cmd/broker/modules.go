package main

import (
	"fmt"

	ac "github.com/Azure/open-service-broker-azure/pkg/azure/aci"
	"github.com/Azure/open-service-broker-azure/pkg/azure/arm"
	cd "github.com/Azure/open-service-broker-azure/pkg/azure/cosmosdb"
	eh "github.com/Azure/open-service-broker-azure/pkg/azure/eventhub"
	hd "github.com/Azure/azure-service-broker/pkg/azure/hdinsight"
	kv "github.com/Azure/open-service-broker-azure/pkg/azure/keyvault"
	ss "github.com/Azure/open-service-broker-azure/pkg/azure/mssql"
	mg "github.com/Azure/open-service-broker-azure/pkg/azure/mysql"
	pg "github.com/Azure/open-service-broker-azure/pkg/azure/postgresql"
	rc "github.com/Azure/open-service-broker-azure/pkg/azure/rediscache"
	se "github.com/Azure/open-service-broker-azure/pkg/azure/search"
	sb "github.com/Azure/open-service-broker-azure/pkg/azure/servicebus"
	sa "github.com/Azure/open-service-broker-azure/pkg/azure/storage"
	"github.com/Azure/open-service-broker-azure/pkg/services/mysqldb"
	"github.com/Azure/open-service-broker-azure/pkg/services/sqldb"

	"github.com/Azure/open-service-broker-azure/pkg/service"
	"github.com/Azure/open-service-broker-azure/pkg/services/aci"
	"github.com/Azure/open-service-broker-azure/pkg/services/cosmosdb"
	"github.com/Azure/open-service-broker-azure/pkg/services/eventhubs"
	"github.com/Azure/azure-service-broker/pkg/services/hdinsight"
	"github.com/Azure/open-service-broker-azure/pkg/services/keyvault"
	"github.com/Azure/open-service-broker-azure/pkg/services/postgresqldb"
	"github.com/Azure/open-service-broker-azure/pkg/services/rediscache"
	"github.com/Azure/open-service-broker-azure/pkg/services/search"
	"github.com/Azure/open-service-broker-azure/pkg/services/servicebus"
	"github.com/Azure/open-service-broker-azure/pkg/services/storage"
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
	serviceBusManager, err := sb.NewManager()
	if err != nil {
		return fmt.Errorf("error initializing service bus manager: %s", err)
	}
	eventHubManager, err := eh.NewManager()
	if err != nil {
		return fmt.Errorf("error initializing event hub manager: %s", err)
	}
	keyvaultManager, err := kv.NewManager()
	if err != nil {
		return fmt.Errorf("error initializing keyvault manager: %s", err)
	}
	msSQLManager, err := ss.NewManager()
	if err != nil {
		return fmt.Errorf("error initializing mssql manager: %s", err)
	}
	msSQLConfig, err := sqldb.GetConfig()
	if err != nil {
		return fmt.Errorf("error parsing mssql configuration: %s", err)
	}
	cosmosDBManager, err := cd.NewManager()
	if err != nil {
		return fmt.Errorf("error initializing cosmosdb manager: %s", err)
	}
	storageManager, err := sa.NewManager()
	if err != nil {
		return fmt.Errorf("error initializing storage manager: %s", err)
	}
	searchManager, err := se.NewManager()
	if err != nil {
		return fmt.Errorf("error initializing search manager: %s", err)
	}
	aciManager, err := ac.NewManager()
	if err != nil {
		return fmt.Errorf("error initializing aci manager: %s", err)
	}
	hdinsightManager, err := hd.NewManager()
	if err != nil {
		return fmt.Errorf("error initializing hdinsight manager: %s", err)
	}

	modules = []service.Module{
		postgresqldb.New(armDeployer, postgreSQLManager),
		rediscache.New(armDeployer, redisManager),
		mysqldb.New(armDeployer, mySQLManager),
		servicebus.New(armDeployer, serviceBusManager),
		eventhubs.New(armDeployer, eventHubManager),
		keyvault.New(armDeployer, keyvaultManager),
		sqldb.New(armDeployer, msSQLManager, msSQLConfig),
		cosmosdb.New(armDeployer, cosmosDBManager),
		storage.New(armDeployer, storageManager),
		search.New(armDeployer, searchManager),
		aci.New(armDeployer, aciManager),
		hdinsight.New(armDeployer, hdinsightManager),
	}
	return nil
}

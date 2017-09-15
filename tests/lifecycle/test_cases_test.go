// +build !unit

package lifecycle

import (
	"github.com/Azure/azure-service-broker/pkg/azure/arm"
	cd "github.com/Azure/azure-service-broker/pkg/azure/cosmosdb"
	ss "github.com/Azure/azure-service-broker/pkg/azure/mssql"
	mg "github.com/Azure/azure-service-broker/pkg/azure/mysql"
	pg "github.com/Azure/azure-service-broker/pkg/azure/postgresql"
	rc "github.com/Azure/azure-service-broker/pkg/azure/rediscache"
	"github.com/Azure/azure-service-broker/pkg/services/cosmosdb"
	"github.com/Azure/azure-service-broker/pkg/services/mssql"
	"github.com/Azure/azure-service-broker/pkg/services/mysql"
	"github.com/Azure/azure-service-broker/pkg/services/postgresql"
	"github.com/Azure/azure-service-broker/pkg/services/rediscache"
	uuid "github.com/satori/go.uuid"
)

func getTestCases() ([]moduleLifecycleTestCase, error) {
	armDeployer, err := arm.NewDeployer()
	if err != nil {
		return nil, err
	}
	postgreSQLManager, err := pg.NewManager()
	if err != nil {
		return nil, err
	}
	mySQLManager, err := mg.NewManager()
	if err != nil {
		return nil, err
	}
	redisManager, err := rc.NewManager()
	if err != nil {
		return nil, err
	}
	msSQLManager, err := ss.NewManager()
	if err != nil {
		return nil, err
	}
	msSQLConfig, err := mssql.GetConfig()
	if err != nil {
		return nil, err
	}
	cosmosdbManager, err := cd.NewManager()
	if err != nil {
		return nil, err
	}
	return []moduleLifecycleTestCase{
		{ // PostgreSQL
			module:    postgresql.New(armDeployer, postgreSQLManager),
			serviceID: "b43b4bba-5741-4d98-a10b-17dc5cee0175",
			planID:    "b2ed210f-6a10-4593-a6c4-964e6b6fad62",
			provisioningParameters: &postgresql.ProvisioningParameters{
				Location:      "southcentralus",
				ResourceGroup: newTestResourceGroupName(),
			},
			bindingParameters: &postgresql.BindingParameters{},
		},
		{ // MySQL
			module:    mysql.New(armDeployer, mySQLManager),
			serviceID: "997b8372-8dac-40ac-ae65-758b4a5075a5",
			planID:    "427559f1-bf2a-45d3-8844-32374a3e58aa",
			provisioningParameters: &mysql.ProvisioningParameters{
				Location:      "southcentralus",
				ResourceGroup: newTestResourceGroupName(),
			},
			bindingParameters: &mysql.BindingParameters{},
		},
		{ // Redis Cache
			module:    rediscache.New(armDeployer, redisManager),
			serviceID: "0346088a-d4b2-4478-aa32-f18e295ec1d9",
			planID:    "362b3d1b-5b57-4289-80ad-4a15a760c29c",
			provisioningParameters: &rediscache.ProvisioningParameters{
				Location:      "southcentralus",
				ResourceGroup: newTestResourceGroupName(),
			},
			bindingParameters: &rediscache.BindingParameters{},
		},
		{ // MS SQL
			module:    mssql.New(armDeployer, msSQLManager, msSQLConfig),
			serviceID: "fb9bc99e-0aa9-11e6-8a8a-000d3a002ed5",
			planID:    "3819fdfa-0aaa-11e6-86f4-000d3a002ed5",
			provisioningParameters: &mssql.ProvisioningParameters{
				Location:      "southcentralus",
				ResourceGroup: newTestResourceGroupName(),
			},
			bindingParameters: &mssql.BindingParameters{},
		},
		{ // DocumentDB
			module:    cosmosdb.New(armDeployer, cosmosdbManager),
			serviceID: "6330de6f-a561-43ea-a15e-b99f44d183e6",
			planID:    "71168d1a-c704-49ff-8c79-214dd3d6f8eb",
			provisioningParameters: &cosmosdb.ProvisioningParameters{
				Location:      "eastus",
				ResourceGroup: newTestResourceGroupName(),
			},
			bindingParameters: &cosmosdb.BindingParameters{},
		},
		{ // MongoDB
			module:    cosmosdb.New(armDeployer, cosmosdbManager),
			serviceID: "8797a079-5346-4e84-8018-b7d5ea5c0e3a",
			planID:    "86fdda05-78d7-4026-a443-1325928e7b02",
			provisioningParameters: &cosmosdb.ProvisioningParameters{
				Location:      "eastus",
				ResourceGroup: newTestResourceGroupName(),
			},
			bindingParameters: &cosmosdb.BindingParameters{},
		},
	}, nil
}

func newTestResourceGroupName() string {
	return "test-" + uuid.NewV4().String()
}

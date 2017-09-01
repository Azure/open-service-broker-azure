// +build !unit

package lifecycle

import (
	"github.com/Azure/azure-service-broker/pkg/azure/arm"
	mg "github.com/Azure/azure-service-broker/pkg/azure/mysql"
	pg "github.com/Azure/azure-service-broker/pkg/azure/postgresql"
	rc "github.com/Azure/azure-service-broker/pkg/azure/rediscache"
	"github.com/Azure/azure-service-broker/pkg/services/mysql"
	"github.com/Azure/azure-service-broker/pkg/services/postgresql"
	"github.com/Azure/azure-service-broker/pkg/services/rediscache"
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
	return []moduleLifecycleTestCase{
		{
			module:    postgresql.New(armDeployer, postgreSQLManager),
			serviceID: "b43b4bba-5741-4d98-a10b-17dc5cee0175",
			planID:    "b2ed210f-6a10-4593-a6c4-964e6b6fad62",
			provisioningParameters: &postgresql.ProvisioningParameters{
				Location: "southcentralus",
			},
			bindingParameters: &postgresql.BindingParameters{},
		},
		{
			module:    mysql.New(armDeployer, mySQLManager),
			serviceID: "997b8372-8dac-40ac-ae65-758b4a5075a5",
			planID:    "427559f1-bf2a-45d3-8844-32374a3e58aa",
			provisioningParameters: &mysql.ProvisioningParameters{
				Location: "southcentralus",
			},
			bindingParameters: &mysql.BindingParameters{},
		},
		{
			module:    rediscache.New(armDeployer, redisManager),
			serviceID: "0346088a-d4b2-4478-aa32-f18e295ec1d9",
			planID:    "362b3d1b-5b57-4289-80ad-4a15a760c29c",
			provisioningParameters: &rediscache.ProvisioningParameters{
				Location: "southcentralus",
			},
			bindingParameters: &rediscache.BindingParameters{},
		},
	}, nil
}

// +build !unit

package lifecycle

import (
	"github.com/Azure/open-service-broker-azure/pkg/azure/arm"
	mg "github.com/Azure/open-service-broker-azure/pkg/azure/mysql"
	"github.com/Azure/open-service-broker-azure/pkg/service"
	"github.com/Azure/open-service-broker-azure/pkg/services/mysql"
)

func getMysqlCases(
	armDeployer arm.Deployer,
	resourceGroup string,
) ([]moduleLifecycleTestCase, error) {
	mySQLManager, err := mg.NewManager()
	if err != nil {
		return nil, err
	}

	return []moduleLifecycleTestCase{
		{
			module:    mysql.New(armDeployer, mySQLManager),
			serviceID: "997b8372-8dac-40ac-ae65-758b4a5075a5",
			planID:    "427559f1-bf2a-45d3-8844-32374a3e58aa",
			standardProvisioningContext: service.StandardProvisioningContext{
				Location: "southcentralus",
			},
			provisioningParameters: &mysql.ProvisioningParameters{},
			bindingParameters:      &mysql.BindingParameters{},
		},
	}, nil
}

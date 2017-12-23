// +build !unit

package lifecycle

import (
	"github.com/Azure/open-service-broker-azure/pkg/azure/arm"
	mg "github.com/Azure/open-service-broker-azure/pkg/azure/mysql"
	"github.com/Azure/open-service-broker-azure/pkg/services/mysqldb"
)

func getMysqlCases(
	armDeployer arm.Deployer,
	resourceGroup string,
) ([]serviceLifecycleTestCase, error) {
	mySQLManager, err := mg.NewManager()
	if err != nil {
		return nil, err
	}

	return []serviceLifecycleTestCase{
		{
			module:    mysqldb.New(armDeployer, mySQLManager),
			serviceID: "997b8372-8dac-40ac-ae65-758b4a5075a5",
			planID:    "427559f1-bf2a-45d3-8844-32374a3e58aa",
			location:  "southcentralus",
			provisioningParameters: &mysqldb.ProvisioningParameters{
				FirewallIPStart: "0.0.0.0",
				FirewallIPEnd:   "255.255.255.255",
			},
			bindingParameters: &mysqldb.BindingParameters{},
		},
	}, nil
}

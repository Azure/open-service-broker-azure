// +build !unit

package lifecycle

import (
	"github.com/Azure/azure-service-broker/pkg/azure/arm"
	ss "github.com/Azure/azure-service-broker/pkg/azure/mssql"
	"github.com/Azure/azure-service-broker/pkg/services/mssql"
)

func getMssqlCases(
	armDeployer arm.Deployer,
) ([]moduleLifecycleTestCase, error) {
	msSQLManager, err := ss.NewManager()
	if err != nil {
		return nil, err
	}
	msSQLConfig, err := mssql.GetConfig()
	if err != nil {
		return nil, err
	}

	return []moduleLifecycleTestCase{
		{
			module:    mssql.New(armDeployer, msSQLManager, msSQLConfig),
			serviceID: "fb9bc99e-0aa9-11e6-8a8a-000d3a002ed5",
			planID:    "3819fdfa-0aaa-11e6-86f4-000d3a002ed5",
			provisioningParameters: &mssql.ProvisioningParameters{
				Location:      "southcentralus",
				ResourceGroup: newTestResourceGroupName(),
			},
			bindingParameters: &mssql.BindingParameters{},
		},
	}, nil
}

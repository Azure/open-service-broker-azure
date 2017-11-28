// +build !unit

package lifecycle

import (
	"github.com/Azure/azure-service-broker/pkg/azure/arm"
	pg "github.com/Azure/azure-service-broker/pkg/azure/postgresql"
	"github.com/Azure/azure-service-broker/pkg/service"
	"github.com/Azure/azure-service-broker/pkg/services/postgresql"
)

func getPostgresqlCases(
	armDeployer arm.Deployer,
	resourceGroup string,
) ([]moduleLifecycleTestCase, error) {
	postgreSQLManager, err := pg.NewManager()
	if err != nil {
		return nil, err
	}

	return []moduleLifecycleTestCase{
		{
			module:    postgresql.New(armDeployer, postgreSQLManager),
			serviceID: "b43b4bba-5741-4d98-a10b-17dc5cee0175",
			planID:    "b2ed210f-6a10-4593-a6c4-964e6b6fad62",
			standardProvisioningContext: service.StandardProvisioningContext{
				Location: "southcentralus",
			},
			provisioningParameters: &postgresql.ProvisioningParameters{
				Extensions: []string{
					"uuid-ossp",
					"postgis",
				},
			},
			bindingParameters: &postgresql.BindingParameters{},
		},
	}, nil
}

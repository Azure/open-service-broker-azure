// +build !unit

package lifecycle

import (
	"github.com/Azure/open-service-broker-azure/pkg/azure/arm"
	pg "github.com/Azure/open-service-broker-azure/pkg/azure/postgresql"
	"github.com/Azure/open-service-broker-azure/pkg/service"
	"github.com/Azure/open-service-broker-azure/pkg/services/postgresqldb"
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
			module:    postgresqldb.New(armDeployer, postgreSQLManager),
			serviceID: "b43b4bba-5741-4d98-a10b-17dc5cee0175",
			planID:    "b2ed210f-6a10-4593-a6c4-964e6b6fad62",
			standardProvisioningContext: service.StandardProvisioningContext{
				Location: "southcentralus",
			},
			provisioningParameters: &postgresqldb.ProvisioningParameters{
				FirewallIPStart: "0.0.0.0",
				FirewallIPEnd:   "255.255.255.0",
				Extensions: []string{
					"uuid-ossp",
					"postgis",
				},
			},
			bindingParameters: &postgresqldb.BindingParameters{},
		},
	}, nil
}

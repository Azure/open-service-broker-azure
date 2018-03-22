package cosmosdb

import (
	"context"

	"github.com/Azure/open-service-broker-azure/pkg/service"
)

func (c *cosmosAccountManager) GetDeprovisioner(
	service.Plan,
) (service.Deprovisioner, error) {
	return service.NewDeprovisioner(
		service.NewDeprovisioningStep("deleteARMDeployment", c.deleteARMDeployment),
		service.NewDeprovisioningStep(
			"deleteCosmosDBServer",
			c.deleteCosmosDBAccount,
		),
	)
}

func (c *cosmosAccountManager) deleteARMDeployment(
	_ context.Context,
	instance service.Instance,
) (service.InstanceDetails, service.SecureInstanceDetails, error) {
	return deleteARMDeployment(c.armDeployer, instance)
}

func (c *cosmosAccountManager) deleteCosmosDBAccount(
	ctx context.Context,
	instance service.Instance,
) (service.InstanceDetails, service.SecureInstanceDetails, error) {
	return deleteCosmosDBAccount(ctx, c.databaseAccountsClient, instance)
}

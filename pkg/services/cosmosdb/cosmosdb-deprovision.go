package cosmosdb

import (
	"context"

	"github.com/Azure/open-service-broker-azure/pkg/service"
)

func (c *cosmosManager) GetDeprovisioner(
	service.Plan,
) (service.Deprovisioner, error) {
	return service.NewDeprovisioner(
		service.NewDeprovisioningStep("deleteARMDeployment", c.deleteARMDeployment),
		service.NewDeprovisioningStep(
			"deleteCosmosDBServer",
			c.deleteCosmosDBServer,
		),
	)
}

func (c *cosmosManager) deleteARMDeployment(
	_ context.Context,
	instance service.Instance,
) (service.InstanceDetails, service.SecureInstanceDetails, error) {
	err := deleteARMDeployment(c.armDeployer, instance)
	if err != nil {
		return nil, nil, err
	}
	return instance.Details, instance.SecureDetails, nil
}

func (c *cosmosManager) deleteCosmosDBServer(
	ctx context.Context,
	instance service.Instance,
) (service.InstanceDetails, service.SecureInstanceDetails, error) {
	err := deleteCosmosDBServer(ctx, c.databaseAccountsClient, instance)
	if err != nil {
		return nil, nil, err
	}
	return instance.Details, instance.SecureDetails, nil
}

package cosmosdb

import (
	"context"

	"github.com/Azure/open-service-broker-azure/pkg/service"
)

func (m *mongoAccountManager) GetDeprovisioner(
	service.Plan,
) (service.Deprovisioner, error) {
	return service.NewDeprovisioner(
		service.NewDeprovisioningStep("deleteARMDeployment", m.deleteARMDeployment),
		service.NewDeprovisioningStep(
			"deleteCosmosDBServer",
			m.deleteCosmosDBAccount,
		),
	)
}

func (m *mongoAccountManager) deleteARMDeployment(
	_ context.Context,
	instance service.Instance,
) (service.InstanceDetails, service.SecureInstanceDetails, error) {
	return deleteARMDeployment(m.armDeployer, instance)
}

func (m *mongoAccountManager) deleteCosmosDBAccount(
	ctx context.Context,
	instance service.Instance,
) (service.InstanceDetails, service.SecureInstanceDetails, error) {
	return deleteCosmosDBAccount(ctx, m.databaseAccountsClient, instance)
}

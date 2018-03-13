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
	err := deleteARMDeployment(m.armDeployer, instance)
	if err != nil {
		return nil, nil, err
	}
	return instance.Details, instance.SecureDetails, nil
}

func (m *mongoAccountManager) deleteCosmosDBAccount(
	ctx context.Context,
	instance service.Instance,
) (service.InstanceDetails, service.SecureInstanceDetails, error) {
	err := deleteCosmosDBAccount(ctx, m.databaseAccountsClient, instance)
	if err != nil {
		return nil, nil, err
	}
	return instance.Details, instance.SecureDetails, nil
}

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
) (service.InstanceDetails, error) {
	if err := deleteARMDeployment(
		m.armDeployer,
		instance.ProvisioningParameters,
		instance.Details.(*cosmosdbInstanceDetails),
	); err != nil {
		return nil, err
	}
	return instance.Details, nil
}

func (m *mongoAccountManager) deleteCosmosDBAccount(
	ctx context.Context,
	instance service.Instance,
) (service.InstanceDetails, error) {
	if err := deleteCosmosDBAccount(
		ctx,
		m.databaseAccountsClient,
		instance.ProvisioningParameters,
		instance.Details.(*cosmosdbInstanceDetails),
	); err != nil {
		return nil, err
	}
	return instance.Details, nil
}

package cosmosdb

import (
	"context"
	"fmt"

	"github.com/Azure/azure-service-broker/pkg/service"
)

func (m *module) GetDeprovisioner(
	string,
	string,
) (service.Deprovisioner, error) {
	return service.NewDeprovisioner(
		service.NewDeprovisioningStep("deleteARMDeployment", m.deleteARMDeployment),
		service.NewDeprovisioningStep(
			"deleteCosmosDBServer",
			m.deleteCosmosDBServer,
		),
	)
}

func (m *module) deleteARMDeployment(
	ctx context.Context, // nolint: unparam
	instanceID string, // nolint: unparam
	serviceID string, // nolint: unparam
	planID string, // nolint: unparam
	provisioningContext service.ProvisioningContext,
) (service.ProvisioningContext, error) {
	pc, ok := provisioningContext.(*cosmosdbProvisioningContext)
	if !ok {
		return nil, fmt.Errorf(
			"error casting provisioningContext as *cosmosdbProvisioningContext",
		)
	}
	if err := m.armDeployer.Delete(
		pc.ARMDeploymentName,
		pc.ResourceGroupName,
	); err != nil {
		return nil, fmt.Errorf("error deleting ARM deployment: %s", err)
	}
	return pc, nil
}

func (m *module) deleteCosmosDBServer(
	ctx context.Context, // nolint: unparam
	instanceID string, // nolint: unparam
	serviceID string, // nolint: unparam
	planID string, // nolint: unparam
	provisioningContext service.ProvisioningContext,
) (service.ProvisioningContext, error) {
	pc, ok := provisioningContext.(*cosmosdbProvisioningContext)
	if !ok {
		return nil, fmt.Errorf(
			"error casting provisioningContext as *cosmosdbProvisioningContext",
		)
	}
	if err := m.cosmosdbManager.DeleteDatabaseAccount(
		pc.DatabaseAccountName,
		pc.ResourceGroupName,
	); err != nil {
		return nil, fmt.Errorf("error deleting cosmosdb server: %s", err)
	}
	return pc, nil
}

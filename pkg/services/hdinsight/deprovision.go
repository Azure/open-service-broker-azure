package hdinsight

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
		service.NewDeprovisioningStep(
			"deleteARMDeployment",
			m.deleteARMDeployment,
		),
		service.NewDeprovisioningStep(
			"deleteCluster",
			m.deleteCluster,
		),
		service.NewDeprovisioningStep(
			"deleteStorageAccount",
			m.deleteStorageAccount,
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
	pc, ok := provisioningContext.(*hdinsightProvisioningContext)
	if !ok {
		return nil, fmt.Errorf(
			"error casting provisioningContext as hdinsightProvisioningContext",
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

func (m *module) deleteCluster(
	ctx context.Context, // nolint: unparam
	instanceID string, // nolint: unparam
	serviceID string, // nolint: unparam
	planID string, // nolint: unparam
	provisioningContext service.ProvisioningContext,
) (service.ProvisioningContext, error) {
	pc, ok := provisioningContext.(*hdinsightProvisioningContext)
	if !ok {
		return nil, fmt.Errorf(
			"error casting provisioningContext as hdinsightProvisioningContext",
		)
	}
	if err := m.hdinsightManager.DeleteCluster(
		pc.ClusterName,
		pc.ResourceGroupName,
	); err != nil {
		return nil, fmt.Errorf("error deleting hdinsight cluster: %s", err)
	}
	return pc, nil
}

func (m *module) deleteStorageAccount(
	ctx context.Context, // nolint: unparam
	instanceID string, // nolint: unparam
	serviceID string, // nolint: unparam
	planID string, // nolint: unparam
	provisioningContext service.ProvisioningContext,
) (service.ProvisioningContext, error) {
	pc, ok := provisioningContext.(*hdinsightProvisioningContext)
	if !ok {
		return nil, fmt.Errorf(
			"error casting provisioningContext as hdinsightProvisioningContext",
		)
	}
	if err := m.hdinsightManager.DeleteStorageAccount(
		pc.StorageAccountName,
		pc.ResourceGroupName,
	); err != nil {
		return nil, fmt.Errorf("error deleting hdinsight storage account: %s", err)
	}
	return pc, nil
}

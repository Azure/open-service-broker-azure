package hdinsight

import (
	"context"
	"fmt"

	"github.com/Azure/open-service-broker-azure/pkg/service"
)

func (s *serviceManager) GetDeprovisioner(
	service.Plan,
) (service.Deprovisioner, error) {
	return service.NewDeprovisioner(
		service.NewDeprovisioningStep(
			"deleteARMDeployment",
			s.deleteARMDeployment,
		),
		service.NewDeprovisioningStep(
			"deleteCluster",
			s.deleteCluster,
		),
		service.NewDeprovisioningStep(
			"deleteStorageAccount",
			s.deleteStorageAccount,
		),
	)
}

func (s *serviceManager) deleteARMDeployment(
	_ context.Context,
	_ string, // instanceID
	_ service.Plan,
	standardProvisioningContext service.StandardProvisioningContext,
	provisioningContext service.ProvisioningContext,
) (service.ProvisioningContext, error) {
	pc, ok := provisioningContext.(*hdinsightProvisioningContext)
	if !ok {
		return nil, fmt.Errorf(
			"error casting provisioningContext as *hdinsightProvisioningContext",
		)
	}
	if err := s.armDeployer.Delete(
		pc.ARMDeploymentName,
		standardProvisioningContext.ResourceGroup,
	); err != nil {
		return nil, fmt.Errorf("error deleting ARM deployment: %s", err)
	}
	return pc, nil
}

func (s *serviceManager) deleteCluster(
	_ context.Context,
	_ string, // instanceID
	_ service.Plan,
	standardProvisioningContext service.StandardProvisioningContext,
	provisioningContext service.ProvisioningContext,
) (service.ProvisioningContext, error) {
	pc, ok := provisioningContext.(*hdinsightProvisioningContext)
	if !ok {
		return nil, fmt.Errorf(
			"error casting provisioningContext as hdinsightProvisioningContext",
		)
	}
	if err := s.hdinsightManager.DeleteCluster(
		pc.ClusterName,
		standardProvisioningContext.ResourceGroup,
	); err != nil {
		return nil, fmt.Errorf("error deleting hdinsight cluster: %s", err)
	}
	return pc, nil
}

func (s *serviceManager) deleteStorageAccount(
	_ context.Context,
	_ string, // instanceID
	_ service.Plan,
	standardProvisioningContext service.StandardProvisioningContext,
	provisioningContext service.ProvisioningContext,
) (service.ProvisioningContext, error) {
	pc, ok := provisioningContext.(*hdinsightProvisioningContext)
	if !ok {
		return nil, fmt.Errorf(
			"error casting provisioningContext as hdinsightProvisioningContext",
		)
	}
	if err := s.hdinsightManager.DeleteStorageAccount(
		pc.StorageAccountName,
		standardProvisioningContext.ResourceGroup,
	); err != nil {
		return nil, fmt.Errorf("error deleting hdinsight storage account: %s", err)
	}
	return pc, nil
}

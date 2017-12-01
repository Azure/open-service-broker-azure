package storage

import (
	"context"
	"fmt"

	"github.com/Azure/azure-service-broker/pkg/service"
)

func (s *serviceManager) GetDeprovisioner(
	service.Plan,
) (service.Deprovisioner, error) {
	return service.NewDeprovisioner(
		service.NewDeprovisioningStep("deleteARMDeployment", s.deleteARMDeployment),
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
	pc, ok := provisioningContext.(*storageProvisioningContext)
	if !ok {
		return nil, fmt.Errorf(
			"error casting provisioningContext as storageProvisioningContext",
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

func (s *serviceManager) deleteStorageAccount(
	_ context.Context,
	_ string, // instanceID
	_ service.Plan,
	standardProvisioningContext service.StandardProvisioningContext,
	provisioningContext service.ProvisioningContext,
) (service.ProvisioningContext, error) {
	pc, ok := provisioningContext.(*storageProvisioningContext)
	if !ok {
		return nil, fmt.Errorf(
			"error casting provisioningContext as storageProvisioningContext",
		)
	}
	if err := s.storageManager.DeleteStorageAccount(
		pc.StorageAccountName,
		standardProvisioningContext.ResourceGroup,
	); err != nil {
		return nil, fmt.Errorf("error deleting storage account: %s", err)
	}
	return pc, nil
}

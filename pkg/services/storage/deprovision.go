package storage

import (
	"context"
	"fmt"

	"github.com/Azure/open-service-broker-azure/pkg/service"
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
	instance service.Instance,
) (service.InstanceDetails, error) {
	dt, ok := instance.Details.(*storageInstanceDetails)
	if !ok {
		return nil, fmt.Errorf(
			"error casting instance.Details as *storageInstanceDetails",
		)
	}
	if err := s.armDeployer.Delete(
		dt.ARMDeploymentName,
		instance.ResourceGroup,
	); err != nil {
		return nil, fmt.Errorf("error deleting ARM deployment: %s", err)
	}
	return dt, nil
}

func (s *serviceManager) deleteStorageAccount(
	_ context.Context,
	instance service.Instance,
) (service.InstanceDetails, error) {
	dt, ok := instance.Details.(*storageInstanceDetails)
	if !ok {
		return nil, fmt.Errorf(
			"error casting instance.Details as *storageInstanceDetails",
		)
	}
	_, err := s.accountsClient.Delete(
		instance.ResourceGroup,
		dt.StorageAccountName,
	)
	if err != nil {
		return nil, fmt.Errorf("error deleting storage account: %s", err)
	}
	return dt, nil
}

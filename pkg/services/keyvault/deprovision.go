package keyvault

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
			"deleteKeyVaultServer",
			s.deleteKeyVaultServer,
		),
	)
}

func (s *serviceManager) deleteARMDeployment(
	_ context.Context,
	instance service.Instance,
	_ service.Plan,
) (service.InstanceDetails, error) {
	dt, ok := instance.Details.(*keyvaultInstanceDetails)
	if !ok {
		return nil, fmt.Errorf(
			"error casting instance.Details as *keyvaultInstanceDetails",
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

func (s *serviceManager) deleteKeyVaultServer(
	_ context.Context,
	instance service.Instance,
	_ service.Plan,
) (service.InstanceDetails, error) {
	dt, ok := instance.Details.(*keyvaultInstanceDetails)
	if !ok {
		return nil, fmt.Errorf(
			"error casting instance.Details as *keyvaultInstanceDetails",
		)
	}
	if err := s.keyvaultManager.DeleteVault(
		dt.KeyVaultName,
		instance.ResourceGroup,
	); err != nil {
		return nil, fmt.Errorf("error deleting key vault: %s", err)
	}
	return dt, nil
}

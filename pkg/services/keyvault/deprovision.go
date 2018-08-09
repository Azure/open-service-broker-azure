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
) (service.InstanceDetails, error) {
	dt := instance.Details.(*instanceDetails)
	if err := s.armDeployer.Delete(
		dt.ARMDeploymentName,
		instance.ProvisioningParameters.GetString("resourceGroup"),
	); err != nil {
		return nil, fmt.Errorf("error deleting ARM deployment: %s", err)
	}
	return instance.Details, nil
}

func (s *serviceManager) deleteKeyVaultServer(
	ctx context.Context,
	instance service.Instance,
) (service.InstanceDetails, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	dt := instance.Details.(*instanceDetails)
	if _, err := s.vaultsClient.Delete(
		ctx,
		instance.ProvisioningParameters.GetString("resourceGroup"),
		dt.KeyVaultName,
	); err != nil {
		return nil, fmt.Errorf("error deleting keyvault: %s", err)
	}
	return instance.Details, nil
}

package aci

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
		service.NewDeprovisioningStep("deleteACIServer", s.deleteACIServer),
	)
}

func (s *serviceManager) deleteARMDeployment(
	_ context.Context,
	_ string, // instanceID
	_ service.Plan,
	standardProvisioningContext service.StandardProvisioningContext,
	provisioningContext service.ProvisioningContext,
) (service.ProvisioningContext, error) {
	pc, ok := provisioningContext.(*aciProvisioningContext)
	if !ok {
		return nil, fmt.Errorf(
			"error casting provisioningContext as *aciProvisioningContext",
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

func (s *serviceManager) deleteACIServer(
	_ context.Context,
	_ string, // instanceID
	_ service.Plan,
	standardProvisioningContext service.StandardProvisioningContext,
	provisioningContext service.ProvisioningContext,
) (service.ProvisioningContext, error) {
	pc, ok := provisioningContext.(*aciProvisioningContext)
	if !ok {
		return nil, fmt.Errorf(
			"error casting provisioningContext as *aciProvisioningContext",
		)
	}
	if err := s.aciManager.DeleteACI(
		pc.ContainerName,
		standardProvisioningContext.ResourceGroup,
	); err != nil {
		return nil, fmt.Errorf("error deleting key vault: %s", err)
	}
	return pc, nil
}

package aci

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
		service.NewDeprovisioningStep("deleteACIServer", s.deleteACIServer),
	)
}

func (s *serviceManager) deleteARMDeployment(
	_ context.Context,
	instance service.Instance,
	_ service.Plan,
	_ service.Instance,
) (service.InstanceDetails, error) {
	dt, ok := instance.Details.(*aciInstanceDetails)
	if !ok {
		return nil, fmt.Errorf(
			"error casting instance.Details as *aciInstanceDetails",
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

func (s *serviceManager) deleteACIServer(
	_ context.Context,
	instance service.Instance,
	_ service.Plan,
	_ service.Instance,
) (service.InstanceDetails, error) {
	dt, ok := instance.Details.(*aciInstanceDetails)
	if !ok {
		return nil, fmt.Errorf(
			"error casting instance.Details as *aciInstanceDetails",
		)
	}
	if err := s.aciManager.DeleteACI(
		dt.ContainerName,
		instance.ResourceGroup,
	); err != nil {
		return nil, fmt.Errorf("error deleting key vault: %s", err)
	}
	return dt, nil
}

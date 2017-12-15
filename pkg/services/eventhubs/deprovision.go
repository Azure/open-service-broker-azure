package eventhubs

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
		service.NewDeprovisioningStep("deleteNamespace", s.deleteNamespace),
	)
}

func (s *serviceManager) deleteARMDeployment(
	_ context.Context,
	instance service.Instance,
	_ service.Plan,
) (service.ProvisioningContext, error) {
	pc, ok := instance.ProvisioningContext.(*eventHubProvisioningContext)
	if !ok {
		return nil, fmt.Errorf(
			"error casting instance.ProvisioningContext as " +
				"*eventHubProvisioningContext",
		)
	}
	if err := s.armDeployer.Delete(
		pc.ARMDeploymentName,
		instance.StandardProvisioningContext.ResourceGroup,
	); err != nil {
		return nil, fmt.Errorf("error deleting ARM deployment: %s", err)
	}
	return pc, nil
}

func (s *serviceManager) deleteNamespace(
	_ context.Context,
	instance service.Instance,
	_ service.Plan,
) (service.ProvisioningContext, error) {
	pc, ok := instance.ProvisioningContext.(*eventHubProvisioningContext)
	if !ok {
		return nil, fmt.Errorf(
			"error casting instance.ProvisioningContext as " +
				"*eventHubProvisioningContext",
		)
	}
	if err := s.eventHubManager.DeleteNamespace(
		instance.StandardProvisioningContext.ResourceGroup,
		pc.EventHubNamespace,
	); err != nil {
		return nil, fmt.Errorf("error deleting event hub namespace: %s", err)
	}
	return pc, nil
}

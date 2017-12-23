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
) (service.InstanceDetails, error) {
	dt, ok := instance.Details.(*eventHubInstanceDetails)
	if !ok {
		return nil, fmt.Errorf(
			"error casting instance.Details as *eventHubInstanceDetails",
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

func (s *serviceManager) deleteNamespace(
	_ context.Context,
	instance service.Instance,
	_ service.Plan,
) (service.InstanceDetails, error) {
	dt, ok := instance.Details.(*eventHubInstanceDetails)
	if !ok {
		return nil, fmt.Errorf(
			"error casting instance.Details as *eventHubInstanceDetails",
		)
	}
	if err := s.eventHubManager.DeleteNamespace(
		instance.ResourceGroup,
		dt.EventHubNamespace,
	); err != nil {
		return nil, fmt.Errorf("error deleting event hub namespace: %s", err)
	}
	return dt, nil
}

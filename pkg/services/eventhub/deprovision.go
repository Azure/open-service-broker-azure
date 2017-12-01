package eventhub

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
		service.NewDeprovisioningStep("deleteNamespace", s.deleteNamespace),
	)
}

func (s *serviceManager) deleteARMDeployment(
	_ context.Context,
	_ string, // instanceID
	_ service.Plan,
	standardProvisioningContext service.StandardProvisioningContext,
	provisioningContext service.ProvisioningContext,
) (service.ProvisioningContext, error) {
	pc, ok := provisioningContext.(*eventHubProvisioningContext)
	if !ok {
		return nil, fmt.Errorf(
			"error casting provisioningContext as *eventHubProvisioningContext",
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

func (s *serviceManager) deleteNamespace(
	_ context.Context,
	_ string, // instanceID
	_ service.Plan,
	standardProvisioningContext service.StandardProvisioningContext,
	provisioningContext service.ProvisioningContext,
) (service.ProvisioningContext, error) {
	pc, ok := provisioningContext.(*eventHubProvisioningContext)
	if !ok {
		return nil, fmt.Errorf(
			"error casting provisioningContext as *eventHubProvisioningContext",
		)
	}
	if err := s.eventHubManager.DeleteNamespace(
		standardProvisioningContext.ResourceGroup,
		pc.EventHubNamespace,
	); err != nil {
		return nil, fmt.Errorf("error deleting event hub namespace: %s", err)
	}
	return pc, nil
}

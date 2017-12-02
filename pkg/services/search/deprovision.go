package search

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
		service.NewDeprovisioningStep("deleteAzureSearch", s.deleteAzureSearch),
	)
}

func (s *serviceManager) deleteARMDeployment(
	_ context.Context,
	_ string, // instanceID
	_ service.Plan,
	standardProvisioningContext service.StandardProvisioningContext,
	provisioningContext service.ProvisioningContext,
) (service.ProvisioningContext, error) {
	pc, ok := provisioningContext.(*searchProvisioningContext)
	if !ok {
		return nil, fmt.Errorf(
			"error casting provisioningContext as searchProvisioningContext",
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

func (s *serviceManager) deleteAzureSearch(
	_ context.Context,
	_ string, // instanceID
	_ service.Plan,
	standardProvisioningContext service.StandardProvisioningContext,
	provisioningContext service.ProvisioningContext,
) (service.ProvisioningContext, error) {
	pc, ok := provisioningContext.(*searchProvisioningContext)
	if !ok {
		return nil, fmt.Errorf(
			"error casting provisioningContext as searchProvisioningContext",
		)
	}
	if err := s.searchManager.DeleteServer(
		pc.ServiceName,
		standardProvisioningContext.ResourceGroup,
	); err != nil {
		return nil, fmt.Errorf("error deleting Azure Search: %s", err)
	}
	return pc, nil
}

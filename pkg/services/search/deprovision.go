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
	instance service.Instance,
) (service.InstanceDetails, error) {
	dt, ok := instance.Details.(*searchInstanceDetails)
	if !ok {
		return nil, fmt.Errorf(
			"error casting instance.Details as *searchInstanceDetails",
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

func (s *serviceManager) deleteAzureSearch(
	ctx context.Context,
	instance service.Instance,
) (service.InstanceDetails, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	dt, ok := instance.Details.(*searchInstanceDetails)
	if !ok {
		return nil, fmt.Errorf(
			"error casting instance.Details as *searchInstanceDetails",
		)
	}
	if _, err := s.servicesClient.Delete(
		ctx,
		instance.ResourceGroup,
		dt.ServiceName,
		nil,
	); err != nil {
		return nil, fmt.Errorf("error deleting Azure Search: %s", err)
	}
	return dt, nil
}

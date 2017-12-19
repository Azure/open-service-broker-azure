package cosmosdb

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
			"deleteCosmosDBServer",
			s.deleteCosmosDBServer,
		),
	)
}

func (s *serviceManager) deleteARMDeployment(
	_ context.Context,
	instance service.Instance,
	_ service.Plan,
	_ service.Instance, // Reference instance
) (service.ProvisioningContext, error) {
	pc, ok := instance.ProvisioningContext.(*cosmosdbProvisioningContext)
	if !ok {
		return nil, fmt.Errorf(
			"error casting instance.ProvisioningContext as " +
				"*cosmosdbProvisioningContext",
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

func (s *serviceManager) deleteCosmosDBServer(
	_ context.Context,
	instance service.Instance,
	_ service.Plan,
	_ service.Instance, // Reference instance
) (service.ProvisioningContext, error) {
	pc, ok := instance.ProvisioningContext.(*cosmosdbProvisioningContext)
	if !ok {
		return nil, fmt.Errorf(
			"error casting instance.ProvisioningContext as " +
				"*cosmosdbProvisioningContext",
		)
	}
	if err := s.cosmosdbManager.DeleteDatabaseAccount(
		pc.DatabaseAccountName,
		instance.StandardProvisioningContext.ResourceGroup,
	); err != nil {
		return nil, fmt.Errorf("error deleting cosmosdb server: %s", err)
	}
	return pc, nil
}

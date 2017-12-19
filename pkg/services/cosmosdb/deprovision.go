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
) (service.InstanceDetails, error) {
	dt, ok := instance.Details.(*cosmosdbInstanceDetails)
	if !ok {
		return nil, fmt.Errorf(
			"error casting instance.Details as *cosmosdbInstanceDetails",
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

func (s *serviceManager) deleteCosmosDBServer(
	_ context.Context,
	instance service.Instance,
	_ service.Plan,
) (service.InstanceDetails, error) {
	dt, ok := instance.Details.(*cosmosdbInstanceDetails)
	if !ok {
		return nil, fmt.Errorf(
			"error casting instance.Details as *cosmosdbInstanceDetails",
		)
	}
	if err := s.cosmosdbManager.DeleteDatabaseAccount(
		dt.DatabaseAccountName,
		instance.ResourceGroup,
	); err != nil {
		return nil, fmt.Errorf("error deleting cosmosdb server: %s", err)
	}
	return dt, nil
}

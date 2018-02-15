package cosmosdb

import (
	"context"
	"fmt"
	"strings"

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
	ctx context.Context,
	instance service.Instance,
) (service.InstanceDetails, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	dt, ok := instance.Details.(*cosmosdbInstanceDetails)
	if !ok {
		return nil, fmt.Errorf(
			"error casting instance.Details as *cosmosdbInstanceDetails",
		)
	}
	result, err := s.databaseAccountsClient.Delete(
		ctx,
		instance.ResourceGroup,
		dt.DatabaseAccountName,
	)
	if err != nil {
		return dt, fmt.Errorf("error deleting cosmosdb server: %s", err)
	}
	if err := result.WaitForCompletion(
		ctx,
		s.databaseAccountsClient.Client,
	); err != nil {
		// Workaround for https://github.com/Azure/azure-sdk-for-go/issues/759
		if strings.Contains(err.Error(), "StatusCode=404") {
			return dt, nil
		}
		return dt, fmt.Errorf("error deleting cosmosdb server: %s", err)
	}
	return dt, nil
}

package eventhubs

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
		service.NewDeprovisioningStep("deleteNamespace", s.deleteNamespace),
	)
}

func (s *serviceManager) deleteARMDeployment(
	_ context.Context,
	instance service.Instance,
) (service.InstanceDetails, error) {
	dt := instance.Details.(*instanceDetails)
	if err := s.armDeployer.Delete(
		dt.ARMDeploymentName,
		instance.ProvisioningParameters.GetString("resourceGroup"),
	); err != nil {
		return nil, fmt.Errorf("error deleting ARM deployment: %s", err)
	}
	return instance.Details, nil
}

func (s *serviceManager) deleteNamespace(
	ctx context.Context,
	instance service.Instance,
) (service.InstanceDetails, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	dt := instance.Details.(*instanceDetails)
	result, err := s.namespacesClient.Delete(
		ctx,
		instance.ProvisioningParameters.GetString("resourceGroup"),
		dt.EventHubNamespace,
	)
	if err != nil {
		return nil, fmt.Errorf("error deleting event hub namespace: %s", err)
	}
	if err := result.WaitForCompletionRef(
		ctx,
		s.namespacesClient.Client,
	); err != nil {
		// Workaround for https://github.com/Azure/azure-sdk-for-go/issues/759
		if !strings.Contains(err.Error(), "StatusCode=404") {
			return nil, fmt.Errorf("error deleting event hub namespace: %s", err)
		}
	}
	return instance.Details, nil
}

package eventhub

import (
	"context"
	"fmt"

	"github.com/Azure/azure-service-broker/pkg/service"
)

func (m *module) GetDeprovisioner(
	string,
	string,
) (service.Deprovisioner, error) {
	return service.NewDeprovisioner(
		service.NewDeprovisioningStep("deleteARMDeployment", m.deleteARMDeployment),
		service.NewDeprovisioningStep("deleteNamespace", m.deleteNamespace),
	)
}

func (m *module) deleteARMDeployment(
	_ context.Context,
	_ string, // instanceID
	_ string, // serviceID
	_ string, // planID
	standardProvisioningContext service.StandardProvisioningContext,
	provisioningContext service.ProvisioningContext,
) (service.ProvisioningContext, error) {
	pc, ok := provisioningContext.(*eventHubProvisioningContext)
	if !ok {
		return nil, fmt.Errorf(
			"error casting provisioningContext as *eventHubProvisioningContext",
		)
	}
	if err := m.armDeployer.Delete(
		pc.ARMDeploymentName,
		standardProvisioningContext.ResourceGroup,
	); err != nil {
		return nil, fmt.Errorf("error deleting ARM deployment: %s", err)
	}
	return pc, nil
}

func (m *module) deleteNamespace(
	_ context.Context,
	_ string, // instanceID
	_ string, // serviceID
	_ string, // planID
	standardProvisioningContext service.StandardProvisioningContext,
	provisioningContext service.ProvisioningContext,
) (service.ProvisioningContext, error) {
	pc, ok := provisioningContext.(*eventHubProvisioningContext)
	if !ok {
		return nil, fmt.Errorf(
			"error casting provisioningContext as *eventHubProvisioningContext",
		)
	}
	if err := m.eventHubManager.DeleteNamespace(
		standardProvisioningContext.ResourceGroup,
		pc.EventHubNamespace,
	); err != nil {
		return nil, fmt.Errorf("error deleting event hub namespace: %s", err)
	}
	return pc, nil
}

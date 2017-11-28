package servicebus

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
	pc, ok := provisioningContext.(*serviceBusProvisioningContext)
	if !ok {
		return nil, fmt.Errorf(
			"error casting provisioningContext as *serviceBusProvisioningContext",
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
	pc, ok := provisioningContext.(*serviceBusProvisioningContext)
	if !ok {
		return nil, fmt.Errorf(
			"error casting provisioningContext as *serviceBusProvisioningContext",
		)
	}
	if err := m.serviceBusManager.DeleteNamespace(
		pc.ServiceBusNamespaceName,
		standardProvisioningContext.ResourceGroup,
	); err != nil {
		return nil, fmt.Errorf("error deleting service bus namespace: %s", err)
	}
	return pc, nil
}

package servicebus

import (
	"context"
	"fmt"
	"strings"

	"github.com/Azure/open-service-broker-azure/pkg/service"
)

func (nm *namespaceManager) GetDeprovisioner(
	service.Plan,
) (service.Deprovisioner, error) {
	return service.NewDeprovisioner(
		service.NewDeprovisioningStep("deleteARMDeployment", nm.deleteARMDeployment),
		service.NewDeprovisioningStep("deleteNamespace", nm.deleteNamespace),
	)
}

func (nm *namespaceManager) deleteARMDeployment(
	_ context.Context,
	instance service.Instance,
) (service.InstanceDetails, error) {
	dt := instance.Details.(*namespaceInstanceDetails)
	if err := nm.armDeployer.Delete(
		dt.ARMDeploymentName,
		instance.ProvisioningParameters.GetString("resourceGroup"),
	); err != nil {
		return nil, fmt.Errorf("error deleting ARM deployment: %s", err)
	}
	return instance.Details, nil
}

func (nm *namespaceManager) deleteNamespace(
	ctx context.Context,
	instance service.Instance,
) (service.InstanceDetails, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	dt := instance.Details.(*namespaceInstanceDetails)
	result, err := nm.namespacesClient.Delete(
		ctx,
		instance.ProvisioningParameters.GetString("resourceGroup"),
		dt.ServiceBusNamespaceName,
	)
	if err != nil {
		return nil, fmt.Errorf("error deleting service bus namespace: %s", err)
	}
	if err := result.WaitForCompletion(
		ctx,
		nm.namespacesClient.Client,
	); err != nil {
		// Workaround for https://github.com/Azure/azure-sdk-for-go/issues/759
		if !strings.Contains(err.Error(), "StatusCode=404") {
			return nil, fmt.Errorf("error deleting service bus namespace: %s", err)
		}
	}
	return instance.Details, nil
}

package iothub

import (
	"context"
	"fmt"

	"github.com/Azure/open-service-broker-azure/pkg/service"
)

func (i *iotHubManager) GetDeprovisioner(
	service.Plan,
) (service.Deprovisioner, error) {
	return service.NewDeprovisioner(
		service.NewDeprovisioningStep("deleteARMDeployment", i.deleteARMDeployment),
		service.NewDeprovisioningStep(
			"deleteIoTHub",
			i.deleteIoTHub,
		),
	)
}

func (i *iotHubManager) deleteARMDeployment(
	_ context.Context,
	instance service.Instance,
) (service.InstanceDetails, error) {
	dt := instance.Details.(*instanceDetails)

	if err := i.armDeployer.Delete(
		dt.ARMDeploymentName,
		instance.ProvisioningParameters.GetString("resourceGroup"),
	); err != nil {
		return nil, fmt.Errorf("error deleting ARM deployment: %s", err)
	}
	return instance.Details, nil
}

func (i *iotHubManager) deleteIoTHub(
	ctx context.Context,
	instance service.Instance,
) (service.InstanceDetails, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	dt := instance.Details.(*instanceDetails)

	_, err := i.iotHubClient.Delete(
		ctx,
		instance.ProvisioningParameters.GetString("resourceGroup"),
		dt.IoTHubName,
	)
	if err != nil {
		return nil, fmt.Errorf("error deleting iot hub %s", err)
	}
	return dt, nil
}

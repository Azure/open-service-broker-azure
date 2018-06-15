package mssql

import (
	"context"
	"fmt"

	"github.com/Azure/open-service-broker-azure/pkg/service"
)

func (d *dbmsManager) GetDeprovisioner(
	service.Plan,
) (service.Deprovisioner, error) {
	return service.NewDeprovisioner(
		service.NewDeprovisioningStep("deleteARMDeployment", d.deleteARMDeployment),
		service.NewDeprovisioningStep(
			"deleteMsSQLServer",
			d.deleteMsSQLServer,
		),
	)
}

func (d *dbmsManager) deleteARMDeployment(
	_ context.Context,
	instance service.Instance,
) (service.InstanceDetails, error) {
	dt := instance.Details.(*dbmsInstanceDetails)
	err := d.armDeployer.Delete(
		dt.ARMDeploymentName,
		instance.ProvisioningParameters.GetString("resourceGroup"),
	)
	if err != nil {
		return nil, fmt.Errorf("error deleting ARM deployment: %s", err)
	}
	return instance.Details, nil
}

func (d *dbmsManager) deleteMsSQLServer(
	ctx context.Context,
	instance service.Instance,
) (service.InstanceDetails, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	dt := instance.Details.(*dbmsInstanceDetails)
	result, err := d.serversClient.Delete(
		ctx,
		instance.ProvisioningParameters.GetString("resourceGroup"),
		dt.ServerName,
	)
	if err != nil {
		return nil, fmt.Errorf("error deleting sql server: %s", err)
	}
	if err := result.WaitForCompletion(ctx, d.serversClient.Client); err != nil {
		return nil, fmt.Errorf("error deleting sql server: %s", err)
	}
	return instance.Details, nil
}

package postgresql

import (
	"context"
	"fmt"

	"open-service-broker-azure/pkg/service"
)

func (d *databaseManager) GetDeprovisioner(
	service.Plan,
) (service.Deprovisioner, error) {
	return service.NewDeprovisioner(
		service.NewDeprovisioningStep("deleteARMDeployment", d.deleteARMDeployment),
		service.NewDeprovisioningStep(
			"deletePostgreSQLDatabase",
			d.deletePostgreSQLDatabase,
		),
	)
}

func (d *databaseManager) deleteARMDeployment(
	_ context.Context,
	instance service.Instance,
) (service.InstanceDetails, service.SecureInstanceDetails, error) {
	dt := databaseInstanceDetails{}
	if err := service.GetStructFromMap(instance.Details, &dt); err != nil {
		return nil, nil, err
	}
	if err := d.armDeployer.Delete(
		dt.ARMDeploymentName,
		instance.Parent.ProvisioningParameters.GetString("resourceGroup"),
	); err != nil {
		return nil, nil, fmt.Errorf("error deleting ARM deployment: %s", err)
	}
	return instance.Details, instance.SecureDetails, nil
}

func (d *databaseManager) deletePostgreSQLDatabase(
	ctx context.Context,
	instance service.Instance,
) (service.InstanceDetails, service.SecureInstanceDetails, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	pdt := dbmsInstanceDetails{}
	if err :=
		service.GetStructFromMap(instance.Parent.Details, &pdt); err != nil {
		return nil, nil, err
	}
	dt := databaseInstanceDetails{}
	if err := service.GetStructFromMap(instance.Details, &dt); err != nil {
		return nil, nil, err
	}
	result, err := d.databasesClient.Delete(
		ctx,
		instance.Parent.ProvisioningParameters.GetString("resourceGroup"),
		pdt.ServerName,
		dt.DatabaseName,
	)
	if err != nil {
		return nil, nil, fmt.Errorf("error deleting postgresql database: %s", err)
	}
	if err := result.WaitForCompletion(ctx, d.databasesClient.Client); err != nil {
		return nil, nil, fmt.Errorf("error deleting postgresql database: %s", err)
	}
	return instance.Details, instance.SecureDetails, nil
}

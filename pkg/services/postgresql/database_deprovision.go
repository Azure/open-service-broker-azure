package postgresql

import (
	"context"
	"fmt"

	"github.com/Azure/open-service-broker-azure/pkg/service"
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
) (service.InstanceDetails, error) {
	dt := instance.Details.(*databaseInstanceDetails)
	if err := d.armDeployer.Delete(
		dt.ARMDeploymentName,
		instance.Parent.ProvisioningParameters.GetString("resourceGroup"),
	); err != nil {
		return nil, fmt.Errorf("error deleting ARM deployment: %s", err)
	}
	return instance.Details, nil
}

func (d *databaseManager) deletePostgreSQLDatabase(
	ctx context.Context,
	instance service.Instance,
) (service.InstanceDetails, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	pdt := instance.Parent.Details.(*dbmsInstanceDetails)
	dt := instance.Details.(*databaseInstanceDetails)
	result, err := d.databasesClient.Delete(
		ctx,
		instance.Parent.ProvisioningParameters.GetString("resourceGroup"),
		pdt.ServerName,
		dt.DatabaseName,
	)
	if err != nil {
		return nil, fmt.Errorf("error deleting postgresql database: %s", err)
	}
	if err := result.WaitForCompletionRef(
		ctx,
		d.databasesClient.Client,
	); err != nil {
		return nil, fmt.Errorf("error deleting postgresql database: %s", err)
	}
	return instance.Details, nil
}

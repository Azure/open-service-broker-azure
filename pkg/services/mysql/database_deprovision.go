package mysql

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
		service.NewDeprovisioningStep("deleteMySQLDatabase", d.deleteMySQLDatabase),
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
		instance.Parent.ResourceGroup,
	); err != nil {
		return nil, nil, fmt.Errorf("error deleting ARM deployment: %s", err)
	}
	return instance.Details, instance.SecureDetails, nil
}

func (d *databaseManager) deleteMySQLDatabase(
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
		instance.Parent.ResourceGroup,
		pdt.ServerName,
		dt.DatabaseName,
	)
	if err != nil {
		return nil, nil, fmt.Errorf("error deleting mysql database: %s", err)
	}
	if err := result.WaitForCompletion(ctx, d.databasesClient.Client); err != nil {
		return nil, nil, fmt.Errorf("error deleting mysql database: %s", err)
	}
	return instance.Details, instance.SecureDetails, nil
}

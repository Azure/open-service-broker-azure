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
	dt, ok := instance.Details.(*databaseInstanceDetails)
	if !ok {
		return nil, nil, fmt.Errorf(
			"error casting instance.Details as *mysql.databaseInstanceDetails",
		)
	}
	if err := d.armDeployer.Delete(
		dt.ARMDeploymentName,
		instance.Parent.ResourceGroup,
	); err != nil {
		return nil, nil, fmt.Errorf("error deleting ARM deployment: %s", err)
	}
	return dt, instance.SecureDetails, nil
}

func (d *databaseManager) deleteMySQLDatabase(
	ctx context.Context,
	instance service.Instance,
) (service.InstanceDetails, service.SecureInstanceDetails, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	pdt, ok := instance.Parent.Details.(*dbmsInstanceDetails)
	if !ok {
		return nil, nil, fmt.Errorf(
			"error casting instance.Details as *mysql.dbmsInstanceDetails",
		)
	}
	dt, ok := instance.Details.(*databaseInstanceDetails)
	if !ok {
		return nil, nil, fmt.Errorf(
			"error casting instance.Details as *mysql.databaseInstanceDetails",
		)
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
	return dt, instance.SecureDetails, nil
}

package mysqldb

import (
	"context"
	"fmt"

	"github.com/Azure/open-service-broker-azure/pkg/service"
)

func (d *dbOnlyManager) GetDeprovisioner(
	service.Plan,
) (service.Deprovisioner, error) {
	return service.NewDeprovisioner(
		service.NewDeprovisioningStep("deleteARMDeployment", d.deleteARMDeployment),
		service.NewDeprovisioningStep("deleteMySQLDatabase", d.deleteMySQLDatabase),
	)
}

func (d *dbOnlyManager) deleteARMDeployment(
	_ context.Context,
	instance service.Instance,
) (service.InstanceDetails, service.SecureInstanceDetails, error) {
	dt, ok := instance.Details.(*dbOnlyMysqlInstanceDetails)
	if !ok {
		return nil, nil, fmt.Errorf(
			"error casting instance.Details as *dbOnlyMysqlInstanceDetails",
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

func (d *dbOnlyManager) deleteMySQLDatabase(
	ctx context.Context,
	instance service.Instance,
) (service.InstanceDetails, service.SecureInstanceDetails, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	pdt, ok := instance.Parent.Details.(*dbmsOnlyMysqlInstanceDetails)
	if !ok {
		return nil, nil, fmt.Errorf(
			"error casting instance.Details as *dbOnlyMysqlInstanceDetails",
		)
	}
	dt, ok := instance.Details.(*dbOnlyMysqlInstanceDetails)
	if !ok {
		return nil, nil, fmt.Errorf(
			"error casting instance.Details as *dbOnlyMysqlInstanceDetails",
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

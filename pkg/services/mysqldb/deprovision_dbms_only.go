package mysqldb

import (
	"context"
	"fmt"

	"github.com/Azure/open-service-broker-azure/pkg/service"
)

func (d *dbmsOnlyManager) GetDeprovisioner(
	service.Plan,
) (service.Deprovisioner, error) {
	return service.NewDeprovisioner(
		service.NewDeprovisioningStep("deleteARMDeployment", d.deleteARMDeployment),
		service.NewDeprovisioningStep("deleteMySQLServer", d.deleteMySQLServer),
	)
}

func (d *dbmsOnlyManager) deleteARMDeployment(
	_ context.Context,
	instance service.Instance,
) (service.InstanceDetails, service.SecureInstanceDetails, error) {
	dt, ok := instance.Details.(*dbmsOnlyMysqlInstanceDetails)
	if !ok {
		return nil, nil, fmt.Errorf(
			"error casting instance.Details as *dbmsOnlyMysqlInstanceDetails",
		)
	}
	if err := d.armDeployer.Delete(
		dt.ARMDeploymentName,
		instance.ResourceGroup,
	); err != nil {
		return nil, nil, fmt.Errorf("error deleting ARM deployment: %s", err)
	}
	return dt, instance.SecureDetails, nil
}

func (d *dbmsOnlyManager) deleteMySQLServer(
	ctx context.Context,
	instance service.Instance,
) (service.InstanceDetails, service.SecureInstanceDetails, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	dt, ok := instance.Details.(*dbmsOnlyMysqlInstanceDetails)
	if !ok {
		return nil, nil, fmt.Errorf(
			"error casting instance.Details as *dbmsOnlyMysqlInstanceDetails",
		)
	}
	result, err := d.serversClient.Delete(
		ctx,
		instance.ResourceGroup,
		dt.ServerName,
	)
	if err != nil {
		return nil, nil, fmt.Errorf("error deleting mysql server: %s", err)
	}
	if err := result.WaitForCompletion(ctx, d.serversClient.Client); err != nil {
		return nil, nil, fmt.Errorf("error deleting mysql server: %s", err)
	}
	return dt, instance.SecureDetails, nil
}

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
) (service.InstanceDetails, error) {
	dt, ok := instance.Details.(*dbmsOnlyMysqlInstanceDetails)
	if !ok {
		return nil, fmt.Errorf(
			"error casting instance.Details as *dbmsOnlyMysqlInstanceDetails",
		)
	}
	if err := d.armDeployer.Delete(
		dt.ARMDeploymentName,
		instance.ResourceGroup,
	); err != nil {
		return nil, fmt.Errorf("error deleting ARM deployment: %s", err)
	}
	return dt, nil
}

func (d *dbmsOnlyManager) deleteMySQLServer(
	_ context.Context,
	instance service.Instance,
) (service.InstanceDetails, error) {
	dt, ok := instance.Details.(*dbmsOnlyMysqlInstanceDetails)
	if !ok {
		return nil, fmt.Errorf(
			"error casting instance.Details as *dbmsOnlyMysqlInstanceDetails",
		)
	}
	cancelCh := make(chan struct{})
	_, errChan := d.serversClient.Delete(
		instance.ResourceGroup,
		dt.ServerName,
		cancelCh,
	)
	if err := <-errChan; err != nil {
		return nil, fmt.Errorf("error deleting mysql server: %s", err)
	}
	return dt, nil
}

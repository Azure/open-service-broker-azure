package mysqldb

import (
	"context"
	"fmt"

	"github.com/Azure/open-service-broker-azure/pkg/service"
)

func (a *allInOneManager) GetDeprovisioner(
	service.Plan,
) (service.Deprovisioner, error) {
	return service.NewDeprovisioner(
		service.NewDeprovisioningStep("deleteARMDeployment", a.deleteARMDeployment),
		service.NewDeprovisioningStep("deleteMySQLServer", a.deleteMySQLServer),
	)
}

func (v *vmOnlyManager) GetDeprovisioner(
	service.Plan,
) (service.Deprovisioner, error) {
	return service.NewDeprovisioner(
		service.NewDeprovisioningStep("deleteARMDeployment", v.deleteARMDeployment),
		service.NewDeprovisioningStep("deleteMySQLServer", v.deleteMySQLServer),
	)
}

func (d *dbOnlyManager) GetDeprovisioner(
	service.Plan,
) (service.Deprovisioner, error) {
	return service.NewDeprovisioner(
		service.NewDeprovisioningStep("deleteARMDeployment", d.deleteARMDeployment),
		service.NewDeprovisioningStep("deleteMySQLServer", d.deleteMySQLServer),
	)
}

func (a *allInOneManager) deleteARMDeployment(
	_ context.Context,
	instance service.Instance,
) (service.InstanceDetails, error) {
	dt, ok := instance.Details.(*allInOneMysqlInstanceDetails)
	if !ok {
		return nil, fmt.Errorf(
			"error casting instance.Details as *allInOneMysqlInstanceDetails",
		)
	}
	if err := a.armDeployer.Delete(
		dt.ARMDeploymentName,
		instance.ResourceGroup,
	); err != nil {
		return nil, fmt.Errorf("error deleting ARM deployment: %s", err)
	}
	return dt, nil
}

func (v *vmOnlyManager) deleteARMDeployment(
	_ context.Context,
	instance service.Instance,
) (service.InstanceDetails, error) {
	dt, ok := instance.Details.(*vmOnlyMysqlInstanceDetails)
	if !ok {
		return nil, fmt.Errorf(
			"error casting instance.Details as *vmOnlyMysqlInstanceDetails",
		)
	}
	if err := v.armDeployer.Delete(
		dt.ARMDeploymentName,
		instance.ResourceGroup,
	); err != nil {
		return nil, fmt.Errorf("error deleting ARM deployment: %s", err)
	}
	return dt, nil
}

func (d *dbOnlyManager) deleteARMDeployment(
	_ context.Context,
	instance service.Instance,
) (service.InstanceDetails, error) {
	dt, ok := instance.Details.(*dbOnlyMysqlInstanceDetails)
	if !ok {
		return nil, fmt.Errorf(
			"error casting instance.Details as *dbOnlyMysqlInstanceDetails",
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

func (a *allInOneManager) deleteMySQLServer(
	_ context.Context,
	instance service.Instance,
) (service.InstanceDetails, error) {
	dt, ok := instance.Details.(*allInOneMysqlInstanceDetails)
	if !ok {
		return nil, fmt.Errorf(
			"error casting instance.Details as *allInOneMysqlInstanceDetails",
		)
	}
	cancelCh := make(chan struct{})
	_, errChan := a.serversClient.Delete(
		instance.ResourceGroup,
		dt.ServerName,
		cancelCh,
	)
	if err := <-errChan; err != nil {
		return nil, fmt.Errorf("error deleting mysql server: %s", err)
	}
	return dt, nil
}

func (v *vmOnlyManager) deleteMySQLServer(
	_ context.Context,
	instance service.Instance,
) (service.InstanceDetails, error) {
	dt, ok := instance.Details.(*vmOnlyMysqlInstanceDetails)
	if !ok {
		return nil, fmt.Errorf(
			"error casting instance.Details as *vmOnlyMysqlInstanceDetails",
		)
	}
	cancelCh := make(chan struct{})
	_, errChan := v.serversClient.Delete(
		instance.ResourceGroup,
		dt.ServerName,
		cancelCh,
	)
	if err := <-errChan; err != nil {
		return nil, fmt.Errorf("error deleting mysql server: %s", err)
	}
	return dt, nil
}

func (d *dbOnlyManager) deleteMySQLServer(
	_ context.Context,
	instance service.Instance,
) (service.InstanceDetails, error) {
	pdt, ok := instance.Parent.Details.(*vmOnlyMysqlInstanceDetails)
	if !ok {
		return nil, fmt.Errorf(
			"error casting instance.Details as *dbOnlyMysqlInstanceDetails",
		)
	}
	dt, ok := instance.Details.(*dbOnlyMysqlInstanceDetails)
	if !ok {
		return nil, fmt.Errorf(
			"error casting instance.Details as *dbOnlyMysqlInstanceDetails",
		)
	}
	cancelCh := make(chan struct{})
	_, errChan := d.serversClient.Delete(
		instance.ResourceGroup,
		pdt.ServerName,
		cancelCh,
	)
	if err := <-errChan; err != nil {
		return nil, fmt.Errorf("error deleting mysql server: %s", err)
	}
	return dt, nil
}

package sqldb

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
		service.NewDeprovisioningStep(
			"deleteMsSQLServerOrDatabase",
			a.deleteMsSQLServerOrDatabase,
		),
	)
}

func (v *vmOnlyManager) GetDeprovisioner(
	service.Plan,
) (service.Deprovisioner, error) {
	return service.NewDeprovisioner(
		service.NewDeprovisioningStep("deleteARMDeployment", v.deleteARMDeployment),
	)
}

func (d *dbOnlyManager) GetDeprovisioner(
	service.Plan,
) (service.Deprovisioner, error) {
	return service.NewDeprovisioner(
		service.NewDeprovisioningStep("deleteARMDeployment", d.deleteARMDeployment),
		service.NewDeprovisioningStep(
			"deleteMsSQLServerOrDatabase",
			d.deleteMsSQLServerOrDatabase,
		),
	)
}

func (a *allInOneManager) deleteARMDeployment(
	_ context.Context,
	instance service.Instance,
	_ service.Plan,
	_ service.Instance, // Reference instance
) (service.InstanceDetails, error) {
	dt, ok := instance.Details.(*mssqlAllInOneInstanceDetails)
	if !ok {
		return nil, fmt.Errorf(
			"error casting instance.Details as *mssqlAllInOneInstanceDetails",
		)
	}
	err := a.armDeployer.Delete(
		dt.ARMDeploymentName,
		instance.ResourceGroup,
	)
	if err != nil {
		return nil, fmt.Errorf("error deleting ARM deployment: %s", err)
	}
	return dt, nil
}

func (v *vmOnlyManager) deleteARMDeployment(
	_ context.Context,
	instance service.Instance,
	_ service.Plan,
	_ service.Instance, // Reference instance
) (service.InstanceDetails, error) {
	dt, ok := instance.Details.(*mssqlVMOnlyInstanceDetails)
	if !ok {
		return nil, fmt.Errorf(
			"error casting instance.Details as *mssqlVMOnlyInstanceDetails",
		)
	}
	err := v.armDeployer.Delete(
		dt.ARMDeploymentName,
		instance.ResourceGroup,
	)
	if err != nil {
		return nil, fmt.Errorf("error deleting ARM deployment: %s", err)
	}
	return dt, nil
}

//TODO: Implement DB only logic
func (d *dbOnlyManager) deleteARMDeployment(
	_ context.Context,
	instance service.Instance,
	_ service.Plan,
	_ service.Instance, // Reference instance
) (service.InstanceDetails, error) {
	return instance.Details, nil
}

func (a *allInOneManager) deleteMsSQLServerOrDatabase(
	_ context.Context,
	instance service.Instance,
	_ service.Plan,
	_ service.Instance, // Reference instance
) (service.InstanceDetails, error) {
	dt, ok := instance.Details.(*mssqlAllInOneInstanceDetails)
	if !ok {
		return nil, fmt.Errorf(
			"error casting instance.Details as *mssqlAllInOneInstanceDetails",
		)
	}
	if err := a.mssqlManager.DeleteServer(
		dt.ServerName,
		instance.ResourceGroup,
	); err != nil {
		return dt, fmt.Errorf("error deleting mssql server: %s", err)
	}
	return dt, nil
}

func (d *dbOnlyManager) deleteMsSQLServerOrDatabase(
	_ context.Context,
	instance service.Instance,
	_ service.Plan,
	_ service.Instance, // Reference instance
) (service.InstanceDetails, error) {
	dt, ok := instance.Details.(*mssqlDBOnlyInstanceDetails)
	if !ok {
		return nil, fmt.Errorf(
			"error casting instance.Details as *mssqlDBOnlyInstanceDetails",
		)
	}
	if err := d.mssqlManager.DeleteServer(
		dt.ServerName,
		instance.ResourceGroup,
	); err != nil {
		return dt, fmt.Errorf("error deleting mssql server: %s", err)
	}
	return dt, nil
}

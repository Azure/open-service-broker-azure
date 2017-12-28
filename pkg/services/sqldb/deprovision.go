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
			"deleteMsSQLServer",
			a.deleteMsSQLServer,
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
			"deleteMsSQLDatabase",
			d.deleteMsSQLDatabase,
		),
	)
}

func (a *allInOneManager) deleteARMDeployment(
	_ context.Context,
	instance service.Instance,
	_ service.Plan,
	_ service.Instance, // Reference instance
) (service.InstanceDetails, error) {
	dt, ok := instance.Details.(*mssqlInstanceDetails)
	if !ok {
		return nil, fmt.Errorf(
			"error casting instance.Details as *mssqlInstanceDetails",
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

func (d *dbOnlyManager) deleteARMDeployment(
	_ context.Context,
	instance service.Instance,
	_ service.Plan,
	referenceInstance service.Instance, // Reference instance
) (service.InstanceDetails, error) {
	dt, ok := instance.Details.(*mssqlInstanceDetails)
	if !ok {
		return nil, fmt.Errorf(
			"error casting instance.Details as *mssqlInstanceDetails",
		)
	}
	err := d.armDeployer.Delete(
		dt.ARMDeploymentName,
		referenceInstance.ResourceGroup,
	)
	if err != nil {
		return nil, fmt.Errorf("error deleting ARM deployment: %s", err)
	}
	return dt, nil
}

func (a *allInOneManager) deleteMsSQLServer(
	_ context.Context,
	instance service.Instance,
	_ service.Plan,
	_ service.Instance, // Reference instance
) (service.InstanceDetails, error) {
	dt, ok := instance.Details.(*mssqlInstanceDetails)
	if !ok {
		return nil, fmt.Errorf(
			"error casting instance.Details as *mssqlInstanceDetails",
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

func (d *dbOnlyManager) deleteMsSQLDatabase(
	_ context.Context,
	instance service.Instance,
	_ service.Plan,
	referenceInstance service.Instance, // Reference instance
) (service.InstanceDetails, error) {
	dt, ok := instance.Details.(*mssqlInstanceDetails)
	if !ok {
		return nil, fmt.Errorf(
			"error casting instance.Details as *mssqlInstanceDetails",
		)
	}
	if err := d.mssqlManager.DeleteDatabase(
		dt.ServerName,
		dt.DatabaseName,
		referenceInstance.ResourceGroup,
	); err != nil {
		return dt, fmt.Errorf("error deleting mssql database: %s", err)
	}
	return dt, nil
}

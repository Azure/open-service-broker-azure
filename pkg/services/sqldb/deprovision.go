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
		service.NewDeprovisioningStep(
			"deleteMsSQLServer",
			v.deleteMsSQLServer,
		),
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
) (service.InstanceDetails, error) {
	dt, ok := instance.Details.(*mssqlDBOnlyInstanceDetails)
	if !ok {
		return nil, fmt.Errorf(
			"error casting instance.Details as *mssqlDBOnlyInstanceDetails",
		)
	}
	//Parent should be set by the framework, but return an error if it is not set.
	if instance.Parent == nil {
		return nil, fmt.Errorf("parent instance not set")
	}
	err := d.armDeployer.Delete(
		dt.ARMDeploymentName,
		instance.Parent.ResourceGroup,
	)
	if err != nil {
		return nil, fmt.Errorf("error deleting ARM deployment: %s", err)
	}
	return dt, nil
}

func (a *allInOneManager) deleteMsSQLServer(
	_ context.Context,
	instance service.Instance,
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

func (v *vmOnlyManager) deleteMsSQLServer(
	_ context.Context,
	instance service.Instance,
) (service.InstanceDetails, error) {
	dt, ok := instance.Details.(*mssqlVMOnlyInstanceDetails)
	if !ok {
		return nil, fmt.Errorf(
			"error casting instance.Details as *mssqlInstanceDetails",
		)
	}
	if err := v.mssqlManager.DeleteServer(
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
) (service.InstanceDetails, error) {
	dt, ok := instance.Details.(*mssqlDBOnlyInstanceDetails)
	if !ok {
		return nil, fmt.Errorf(
			"error casting instance.Details as *mssqlDBOnlyInstanceDetails",
		)
	}
	//Parent should be set by the framework, but return an error if it is not set.
	if instance.Parent == nil {
		return nil, fmt.Errorf("parent instance not set")
	}
	pdt, ok := instance.Parent.Details.(*mssqlVMOnlyInstanceDetails)
	if !ok {
		return nil, fmt.Errorf(
			"error casting instance.Parent.Details as " +
				"*mssqlVMOnlyInstanceDetails",
		)
	}

	if err := d.mssqlManager.DeleteDatabase(
		pdt.ServerName,
		dt.DatabaseName,
		instance.Parent.ResourceGroup,
	); err != nil {
		return dt, fmt.Errorf("error deleting mssql database: %s", err)
	}
	return dt, nil
}

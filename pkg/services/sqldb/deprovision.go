package sqldb

import (
	"context"
	"fmt"

	"github.com/Azure/open-service-broker-azure/pkg/service"
)

func (a *allServiceManager) GetDeprovisioner(
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

func (s *vmServiceManager) GetDeprovisioner(
	service.Plan,
) (service.Deprovisioner, error) {
	return service.NewDeprovisioner(
		service.NewDeprovisioningStep("deleteARMDeployment", s.deleteARMDeployment),
		service.NewDeprovisioningStep(
			"deleteMsSQLServerOrDatabase",
			s.deleteMsSQLServerOrDatabase,
		),
	)
}

func (d *dbServiceManager) GetDeprovisioner(
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

func (a *allServiceManager) deleteARMDeployment(
	_ context.Context,
	instance service.Instance,
	_ service.Plan,
	_ service.Instance, // Reference instance
) (service.ProvisioningContext, error) {
	pc, ok := instance.ProvisioningContext.(*mssqlAllInOneProvisioningContext)
	if !ok {
		return nil, fmt.Errorf(
			"error casting instance.ProvisioningContext as " +
				"*mssql.mssqlAllInOneProvisioningContext",
		)
	}
	err := a.armDeployer.Delete(
		pc.ARMDeploymentName,
		instance.StandardProvisioningContext.ResourceGroup,
	)
	if err != nil {
		return nil, fmt.Errorf("error deleting ARM deployment: %s", err)
	}
	return pc, nil
}

func (a *allServiceManager) deleteMsSQLServerOrDatabase(
	_ context.Context,
	instance service.Instance,
	_ service.Plan,
	_ service.Instance, // Reference instance
) (service.ProvisioningContext, error) {
	pc, ok := instance.ProvisioningContext.(*mssqlAllInOneProvisioningContext)
	if !ok {
		return nil, fmt.Errorf(
			"error casting instance.ProvisioningContext as " +
				"*mssql.mssqlAllInOneProvisioningContext",
		)
	}
	if err := a.mssqlManager.DeleteServer(
		pc.ServerName,
		instance.StandardProvisioningContext.ResourceGroup,
	); err != nil {
		return pc, fmt.Errorf("error deleting mssql server: %s", err)
	}
	return pc, nil
}

//TODO: implement db only
func (d *dbServiceManager) deleteARMDeployment(
	_ context.Context,
	instance service.Instance,
	_ service.Plan,
	_ service.Instance, // Reference instance
) (service.ProvisioningContext, error) {
	return instance.ProvisioningContext, nil
}

//TODO: implement db only
func (d *dbServiceManager) deleteMsSQLServerOrDatabase(
	_ context.Context,
	instance service.Instance,
	_ service.Plan,
	_ service.Instance, // Reference instance
) (service.ProvisioningContext, error) {
	return instance.ProvisioningContext, nil
}
func (s *vmServiceManager) deleteARMDeployment(
	_ context.Context,
	instance service.Instance,
	_ service.Plan,
	_ service.Instance, // Reference instance
) (service.ProvisioningContext, error) {
	pc, ok := instance.ProvisioningContext.(*mssqlServerOnlyProvisioningContext)
	if !ok {
		return nil, fmt.Errorf(
			"error casting instance.ProvisioningContext as " +
				"*mssql.mssqlServerOnlyProvisioningContext",
		)
	}
	err := s.armDeployer.Delete(
		pc.ARMDeploymentName,
		instance.StandardProvisioningContext.ResourceGroup,
	)
	if err != nil {
		return nil, fmt.Errorf("error deleting ARM deployment: %s", err)
	}
	return pc, nil
}

func (s *vmServiceManager) deleteMsSQLServerOrDatabase(
	_ context.Context,
	instance service.Instance,
	_ service.Plan,
	_ service.Instance, // Reference instance
) (service.ProvisioningContext, error) {
	pc, ok := instance.ProvisioningContext.(*mssqlServerOnlyProvisioningContext)
	if !ok {
		return nil, fmt.Errorf(
			"error casting instance.ProvisioningContext as " +
				"*mssql.mssqlServerOnlyProvisioningContext",
		)
	}
	if err := s.mssqlManager.DeleteServer(
		pc.ServerName,
		instance.StandardProvisioningContext.ResourceGroup,
	); err != nil {
		return pc, fmt.Errorf("error deleting mssql server: %s", err)
	}
	return pc, nil
}

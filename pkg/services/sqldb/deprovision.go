package sqldb

import (
	"context"
	"fmt"

	"github.com/Azure/open-service-broker-azure/pkg/service"
)

func (s *serviceManager) GetDeprovisioner(
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

func (s *serviceManager) deleteARMDeployment(
	_ context.Context,
	instance service.Instance,
	_ service.Plan,
) (service.ProvisioningContext, error) {
	pc, ok := instance.ProvisioningContext.(*mssqlProvisioningContext)
	if !ok {
		return nil, fmt.Errorf(
			"error casting instance.ProvisioningContext as *mssqlProvisioningContext",
		)
	}
	var err error
	if pc.IsNewServer {
		// new server scenario
		err = s.armDeployer.Delete(
			pc.ARMDeploymentName,
			instance.ResourceGroup,
		)
	} else {
		// exisiting server scenario
		servers := s.mssqlConfig.Servers
		server, ok := servers[pc.ServerName]
		if !ok {
			return nil, fmt.Errorf(
				`can't find serverName "%s" in Azure SQL Server configuration`,
				pc.ServerName,
			)
		}

		err = s.armDeployer.Delete(
			pc.ARMDeploymentName,
			server.ResourceGroupName,
		)
	}
	if err != nil {
		return nil, fmt.Errorf("error deleting ARM deployment: %s", err)
	}
	return pc, nil
}

func (s *serviceManager) deleteMsSQLServerOrDatabase(
	_ context.Context,
	instance service.Instance,
	_ service.Plan,
) (service.ProvisioningContext, error) {
	pc, ok := instance.ProvisioningContext.(*mssqlProvisioningContext)
	if !ok {
		return nil, fmt.Errorf(
			"error casting instance.ProvisioningContext as *mssqlProvisioningContext",
		)
	}

	if pc.IsNewServer {
		// new server scenario
		if err := s.mssqlManager.DeleteServer(
			pc.ServerName,
			instance.ResourceGroup,
		); err != nil {
			return pc, fmt.Errorf("error deleting mssql server: %s", err)
		}
	} else {
		// exisiting server scenario
		servers := s.mssqlConfig.Servers
		server, ok := servers[pc.ServerName]
		if !ok {
			return nil, fmt.Errorf(
				`can't find serverName "%s" in Azure SQL Server configuration`,
				pc.ServerName,
			)
		}

		if err := s.mssqlManager.DeleteDatabase(
			pc.ServerName,
			pc.DatabaseName,
			server.ResourceGroupName,
		); err != nil {
			return pc, fmt.Errorf("error deleting mssql database: %s", err)
		}
	}
	return pc, nil
}

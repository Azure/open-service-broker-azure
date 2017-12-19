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
) (service.InstanceDetails, error) {
	dt, ok := instance.Details.(*mssqlInstanceDetails)
	if !ok {
		return nil, fmt.Errorf(
			"error casting instance.Details as *mssqlInstanceDetails",
		)
	}
	var err error
	if dt.IsNewServer {
		// new server scenario
		err = s.armDeployer.Delete(
			dt.ARMDeploymentName,
			instance.ResourceGroup,
		)
	} else {
		// exisiting server scenario
		servers := s.mssqlConfig.Servers
		server, ok := servers[dt.ServerName]
		if !ok {
			return nil, fmt.Errorf(
				`can't find serverName "%s" in Azure SQL Server configuration`,
				dt.ServerName,
			)
		}

		err = s.armDeployer.Delete(
			dt.ARMDeploymentName,
			server.ResourceGroupName,
		)
	}
	if err != nil {
		return nil, fmt.Errorf("error deleting ARM deployment: %s", err)
	}
	return dt, nil
}

func (s *serviceManager) deleteMsSQLServerOrDatabase(
	_ context.Context,
	instance service.Instance,
	_ service.Plan,
) (service.InstanceDetails, error) {
	dt, ok := instance.Details.(*mssqlInstanceDetails)
	if !ok {
		return nil, fmt.Errorf(
			"error casting instance.Details as *mssqlInstanceDetails",
		)
	}

	if dt.IsNewServer {
		// new server scenario
		if err := s.mssqlManager.DeleteServer(
			dt.ServerName,
			instance.ResourceGroup,
		); err != nil {
			return dt, fmt.Errorf("error deleting mssql server: %s", err)
		}
	} else {
		// exisiting server scenario
		servers := s.mssqlConfig.Servers
		server, ok := servers[dt.ServerName]
		if !ok {
			return nil, fmt.Errorf(
				`can't find serverName "%s" in Azure SQL Server configuration`,
				dt.ServerName,
			)
		}

		if err := s.mssqlManager.DeleteDatabase(
			dt.ServerName,
			dt.DatabaseName,
			server.ResourceGroupName,
		); err != nil {
			return dt, fmt.Errorf("error deleting mssql database: %s", err)
		}
	}
	return dt, nil
}

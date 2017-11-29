package mssql

import (
	"context"
	"fmt"

	"github.com/Azure/azure-service-broker/pkg/service"
)

func (m *module) GetDeprovisioner(
	string,
	string,
) (service.Deprovisioner, error) {
	return service.NewDeprovisioner(
		service.NewDeprovisioningStep(
			"deleteARMDeployment",
			m.deleteARMDeployment,
		),
		service.NewDeprovisioningStep(
			"deleteMsSQLServerOrDatabase",
			m.deleteMsSQLServerOrDatabase,
		),
	)
}

func (m *module) deleteARMDeployment(
	_ context.Context,
	_ string, // instanceID
	_ string, // serviceID
	_ string, // planID
	standardProvisioningContext service.StandardProvisioningContext,
	provisioningContext service.ProvisioningContext,
) (service.ProvisioningContext, error) {
	pc, ok := provisioningContext.(*mssqlProvisioningContext)
	if !ok {
		return nil, fmt.Errorf(
			"error casting provisioningContext as *mssqlProvisioningContext",
		)
	}
	var err error
	if pc.IsNewServer {
		// new server scenario
		err = m.armDeployer.Delete(
			pc.ARMDeploymentName,
			standardProvisioningContext.ResourceGroup,
		)
	} else {
		// exisiting server scenario
		servers := m.mssqlConfig.Servers
		server, ok := servers[pc.ServerName]
		if !ok {
			return nil, fmt.Errorf(
				`can't find serverName "%s" in Azure SQL Server configuration`,
				pc.ServerName,
			)
		}

		err = m.armDeployer.Delete(
			pc.ARMDeploymentName,
			server.ResourceGroupName,
		)
	}
	if err != nil {
		return nil, fmt.Errorf("error deleting ARM deployment: %s", err)
	}
	return pc, nil
}

func (m *module) deleteMsSQLServerOrDatabase(
	_ context.Context,
	_ string, // instanceID
	_ string, // serviceID
	_ string, // planID
	standardProvisioningContext service.StandardProvisioningContext,
	provisioningContext service.ProvisioningContext,
) (service.ProvisioningContext, error) {
	pc, ok := provisioningContext.(*mssqlProvisioningContext)
	if !ok {
		return nil, fmt.Errorf(
			"error casting provisioningContext as *mssqlProvisioningContext",
		)
	}

	if pc.IsNewServer {
		// new server scenario
		if err := m.mssqlManager.DeleteServer(
			pc.ServerName,
			standardProvisioningContext.ResourceGroup,
		); err != nil {
			return pc, fmt.Errorf("error deleting mssql server: %s", err)
		}
	} else {
		// exisiting server scenario
		servers := m.mssqlConfig.Servers
		server, ok := servers[pc.ServerName]
		if !ok {
			return nil, fmt.Errorf(
				`can't find serverName "%s" in Azure SQL Server configuration`,
				pc.ServerName,
			)
		}

		if err := m.mssqlManager.DeleteDatabase(
			pc.ServerName,
			pc.DatabaseName,
			server.ResourceGroupName,
		); err != nil {
			return pc, fmt.Errorf("error deleting mssql database: %s", err)
		}
	}
	return pc, nil
}

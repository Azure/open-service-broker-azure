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
	ctx context.Context, // nolint: unparam
	instanceID string, // nolint: unparam
	serviceID string, // nolint: unparam
	planID string, // nolint: unparam
	provisioningContext service.ProvisioningContext,
) (service.ProvisioningContext, error) {
	pc, ok := provisioningContext.(*mssqlProvisioningContext)
	if !ok {
		return nil, fmt.Errorf(
			"error casting provisioningContext as *mssqlProvisioningContext",
		)
	}
	if err := m.armDeployer.Delete(
		pc.ARMDeploymentName,
		pc.ResourceGroupName,
	); err != nil {
		return nil, fmt.Errorf("error deleting ARM deployment: %s", err)
	}
	return pc, nil
}

func (m *module) deleteMsSQLServerOrDatabase(
	ctx context.Context, // nolint: unparam
	instanceID string, // nolint: unparam
	serviceID string, // nolint: unparam
	planID string, // nolint: unparam
	provisioningContext service.ProvisioningContext,
) (service.ProvisioningContext, error) {
	pc, ok := provisioningContext.(*mssqlProvisioningContext)
	if !ok {
		return nil, fmt.Errorf(
			"error casting provisioningContext as *mssqlProvisioningContext",
		)
	}

	if pc.IsNewServer {
		if err := m.mssqlManager.DeleteServer(
			pc.ServerName,
			pc.ResourceGroupName,
		); err != nil {
			return pc, fmt.Errorf("error deleting mssql server: %s", err)
		}
	} else {
		if err := m.mssqlManager.DeleteDatabase(
			pc.ServerName,
			pc.DatabaseName,
			pc.ResourceGroupName,
		); err != nil {
			return pc, fmt.Errorf("error deleting mssql database: %s", err)
		}
	}
	return pc, nil
}

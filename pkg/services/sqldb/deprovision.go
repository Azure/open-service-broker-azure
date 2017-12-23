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
	_ service.Instance, // Reference instance
) (service.InstanceDetails, error) {
	dt, ok := instance.Details.(*mssqlInstanceDetails)
	if !ok {
		return nil, fmt.Errorf(
			"error casting instance.Details as *mssqlInstanceDetails",
		)
	}
	err := s.armDeployer.Delete(
		dt.ARMDeploymentName,
		instance.ResourceGroup,
	)
	if err != nil {
		return nil, fmt.Errorf("error deleting ARM deployment: %s", err)
	}
	return dt, nil
}

func (s *serviceManager) deleteMsSQLServerOrDatabase(
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
	if err := s.mssqlManager.DeleteServer(
		dt.ServerName,
		instance.ResourceGroup,
	); err != nil {
		return dt, fmt.Errorf("error deleting mssql server: %s", err)
	}
	return dt, nil
}

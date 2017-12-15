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
		service.NewDeprovisioningStep(
			"deleteARMDeployment",
			s.deleteARMDeployment,
		),
		service.NewDeprovisioningStep(
			"deleteMsSQLServerOrDatabase",
			s.deleteMsSQLServerOrDatabase,
		),
	)
}

func (s *serviceManager) deleteARMDeployment(
	_ context.Context,
	_ string, // instanceID
	_ service.Plan,
	standardProvisioningContext service.StandardProvisioningContext,
	provisioningContext service.ProvisioningContext,
) (service.ProvisioningContext, error) {
	pc, ok := provisioningContext.(*mssqlProvisioningContext)
	if !ok {
		return nil, fmt.Errorf(
			"error casting provisioningContext as *mssqlProvisioningContext",
		)
	}
	err := s.armDeployer.Delete(
		pc.ARMDeploymentName,
		standardProvisioningContext.ResourceGroup,
	)
	if err != nil {
		return nil, fmt.Errorf("error deleting ARM deployment: %s", err)
	}
	return pc, nil
}

func (s *serviceManager) deleteMsSQLServerOrDatabase(
	_ context.Context,
	_ string, // instanceID
	_ service.Plan,
	standardProvisioningContext service.StandardProvisioningContext,
	provisioningContext service.ProvisioningContext,
) (service.ProvisioningContext, error) {
	pc, ok := provisioningContext.(*mssqlProvisioningContext)
	if !ok {
		return nil, fmt.Errorf(
			"error casting provisioningContext as *mssqlProvisioningContext",
		)
	}
	if err := s.mssqlManager.DeleteServer(
		pc.ServerName,
		standardProvisioningContext.ResourceGroup,
	); err != nil {
		return pc, fmt.Errorf("error deleting mssql server: %s", err)
	}
	return pc, nil
}

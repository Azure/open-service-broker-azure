package mysqldb

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
		service.NewDeprovisioningStep("deleteMySQLServer", s.deleteMySQLServer),
	)
}

func (s *serviceManager) deleteARMDeployment(
	_ context.Context,
	instance service.Instance,
	_ service.Plan,
	_ service.Instance, // Reference instance
) (service.ProvisioningContext, error) {
	pc, ok := instance.ProvisioningContext.(*mysqlProvisioningContext)
	if !ok {
		return nil, fmt.Errorf(
			"error casting instance.ProvisioningContext as *mysqlProvisioningContext",
		)
	}
	if err := s.armDeployer.Delete(
		pc.ARMDeploymentName,
		instance.StandardProvisioningContext.ResourceGroup,
	); err != nil {
		return nil, fmt.Errorf("error deleting ARM deployment: %s", err)
	}
	return pc, nil
}

func (s *serviceManager) deleteMySQLServer(
	_ context.Context,
	instance service.Instance,
	_ service.Plan,
	_ service.Instance, // Reference instance
) (service.ProvisioningContext, error) {
	pc, ok := instance.ProvisioningContext.(*mysqlProvisioningContext)
	if !ok {
		return nil, fmt.Errorf(
			"error casting instance.ProvisioningContext as *mysqlProvisioningContext",
		)
	}
	if err := s.mysqlManager.DeleteServer(
		pc.ServerName,
		instance.StandardProvisioningContext.ResourceGroup,
	); err != nil {
		return nil, fmt.Errorf("error deleting mysql server: %s", err)
	}
	return pc, nil
}

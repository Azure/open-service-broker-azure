package rediscache

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
		service.NewDeprovisioningStep("deleteRedisServer", s.deleteRedisServer),
	)
}

func (s *serviceManager) deleteARMDeployment(
	_ context.Context,
	instance service.Instance,
	_ service.Plan,
) (service.InstanceDetails, error) {
	dt, ok := instance.Details.(*redisInstanceDetails)
	if !ok {
		return nil, fmt.Errorf(
			"error casting instance.Details as *redisInstanceDetails",
		)
	}
	if err := s.armDeployer.Delete(
		dt.ARMDeploymentName,
		instance.ResourceGroup,
	); err != nil {
		return nil, fmt.Errorf("error deleting ARM deployment: %s", err)
	}
	return dt, nil
}

func (s *serviceManager) deleteRedisServer(
	_ context.Context,
	instance service.Instance,
	_ service.Plan,
) (service.InstanceDetails, error) {
	dt, ok := instance.Details.(*redisInstanceDetails)
	if !ok {
		return nil, fmt.Errorf(
			"error casting instance.Details as *redisInstanceDetails",
		)
	}
	if err := s.redisManager.DeleteServer(
		dt.ServerName,
		instance.ResourceGroup,
	); err != nil {
		return nil, fmt.Errorf("error deleting redis server: %s", err)
	}
	return dt, nil
}

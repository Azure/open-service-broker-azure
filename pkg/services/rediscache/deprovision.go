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
) (service.InstanceDetails, service.SecureInstanceDetails, error) {
	dt, ok := instance.Details.(*redisInstanceDetails)
	if !ok {
		return nil, nil, fmt.Errorf(
			"error casting instance.Details as *redisInstanceDetails",
		)
	}
	if err := s.armDeployer.Delete(
		dt.ARMDeploymentName,
		instance.ResourceGroup,
	); err != nil {
		return nil, nil, fmt.Errorf("error deleting ARM deployment: %s", err)
	}
	return dt, instance.SecureDetails, nil
}

func (s *serviceManager) deleteRedisServer(
	ctx context.Context,
	instance service.Instance,
) (service.InstanceDetails, service.SecureInstanceDetails, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	dt, ok := instance.Details.(*redisInstanceDetails)
	if !ok {
		return nil, nil, fmt.Errorf(
			"error casting instance.Details as *redisInstanceDetails",
		)
	}
	result, err := s.client.Delete(
		ctx,
		instance.ResourceGroup,
		dt.ServerName,
	)
	if err != nil {
		return nil, nil, fmt.Errorf("error deleting redis server: %s", err)
	}
	if err := result.WaitForCompletion(ctx, s.client.Client); err != nil {
		return nil, nil, fmt.Errorf("error deleting redis server: %s", err)
	}
	return dt, instance.SecureDetails, nil
}

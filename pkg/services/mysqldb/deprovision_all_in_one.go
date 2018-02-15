package mysqldb

import (
	"context"
	"fmt"

	"github.com/Azure/open-service-broker-azure/pkg/service"
)

func (a *allInOneManager) GetDeprovisioner(
	service.Plan,
) (service.Deprovisioner, error) {
	return service.NewDeprovisioner(
		service.NewDeprovisioningStep("deleteARMDeployment", a.deleteARMDeployment),
		service.NewDeprovisioningStep("deleteMySQLServer", a.deleteMySQLServer),
	)
}

func (a *allInOneManager) deleteARMDeployment(
	_ context.Context,
	instance service.Instance,
) (service.InstanceDetails, service.SecureInstanceDetails, error) {
	dt, ok := instance.Details.(*allInOneMysqlInstanceDetails)
	if !ok {
		return nil, nil, fmt.Errorf(
			"error casting instance.Details as *allInOneMysqlInstanceDetails",
		)
	}
	if err := a.armDeployer.Delete(
		dt.ARMDeploymentName,
		instance.ResourceGroup,
	); err != nil {
		return nil, nil, fmt.Errorf("error deleting ARM deployment: %s", err)
	}
	return dt, instance.SecureDetails, nil
}

func (a *allInOneManager) deleteMySQLServer(
	ctx context.Context,
	instance service.Instance,
) (service.InstanceDetails, service.SecureInstanceDetails, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	dt, ok := instance.Details.(*allInOneMysqlInstanceDetails)
	if !ok {
		return nil, nil, fmt.Errorf(
			"error casting instance.Details as *allInOneMysqlInstanceDetails",
		)
	}
	result, err := a.serversClient.Delete(
		ctx,
		instance.ResourceGroup,
		dt.ServerName,
	)
	if err != nil {
		return nil, nil, fmt.Errorf("error deleting mysql server: %s", err)
	}
	if err := result.WaitForCompletion(ctx, a.serversClient.Client); err != nil {
		return nil, nil, fmt.Errorf("error deleting mysql server: %s", err)
	}
	return dt, instance.SecureDetails, nil
}

package sqldb

import (
	"context"
	"fmt"

	"github.com/Azure/open-service-broker-azure/pkg/service"
)

func (v *vmOnlyManager) GetDeprovisioner(
	service.Plan,
) (service.Deprovisioner, error) {
	return service.NewDeprovisioner(
		service.NewDeprovisioningStep("deleteARMDeployment", v.deleteARMDeployment),
		service.NewDeprovisioningStep(
			"deleteMsSQLServer",
			v.deleteMsSQLServer,
		),
	)
}

func (v *vmOnlyManager) deleteARMDeployment(
	_ context.Context,
	instance service.Instance,
) (service.InstanceDetails, service.SecureInstanceDetails, error) {
	dt, ok := instance.Details.(*mssqlVMOnlyInstanceDetails)
	if !ok {
		return nil, nil, fmt.Errorf(
			"error casting instance.Details as *mssqlVMOnlyInstanceDetails",
		)
	}
	err := v.armDeployer.Delete(
		dt.ARMDeploymentName,
		instance.ResourceGroup,
	)
	if err != nil {
		return nil, nil, fmt.Errorf("error deleting ARM deployment: %s", err)
	}
	return dt, instance.SecureDetails, nil
}

func (v *vmOnlyManager) deleteMsSQLServer(
	ctx context.Context,
	instance service.Instance,
) (service.InstanceDetails, service.SecureInstanceDetails, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	dt, ok := instance.Details.(*mssqlVMOnlyInstanceDetails)
	if !ok {
		return nil, nil, fmt.Errorf(
			"error casting instance.Details as *mssqlInstanceDetails",
		)
	}
	result, err := v.serversClient.Delete(
		ctx,
		instance.ResourceGroup,
		dt.ServerName,
	)
	if err != nil {
		return nil, nil, fmt.Errorf("error deleting sql server: %s", err)
	}
	if err := result.WaitForCompletion(ctx, v.serversClient.Client); err != nil {
		return nil, nil, fmt.Errorf("error deleting sql server: %s", err)
	}
	return dt, instance.SecureDetails, nil
}

package postgresqldb

import (
	"context"
	"fmt"

	"github.com/Azure/open-service-broker-azure/pkg/service"
)

func (d *dbmsOnlyManager) GetDeprovisioner(
	service.Plan,
) (service.Deprovisioner, error) {
	return service.NewDeprovisioner(
		service.NewDeprovisioningStep("deleteARMDeployment", d.deleteARMDeployment),
		service.NewDeprovisioningStep(
			"deletePostgreSQLServer",
			d.deletePostgreSQLServer,
		),
	)
}

func (d *dbmsOnlyManager) deleteARMDeployment(
	_ context.Context,
	instance service.Instance,
) (service.InstanceDetails, service.SecureInstanceDetails, error) {
	dt, ok := instance.Details.(*dbmsOnlyPostgresqlInstanceDetails)
	if !ok {
		return nil, nil, fmt.Errorf(
			"error casting instance.Details " +
				"as *dbmsOnlyPostgresqlInstanceDetails",
		)
	}
	if err := d.armDeployer.Delete(
		dt.ARMDeploymentName,
		instance.ResourceGroup,
	); err != nil {
		return nil, nil, fmt.Errorf("error deleting ARM deployment: %s", err)
	}
	return dt, instance.SecureDetails, nil
}

func (d *dbmsOnlyManager) deletePostgreSQLServer(
	ctx context.Context,
	instance service.Instance,
) (service.InstanceDetails, service.SecureInstanceDetails, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	dt, ok := instance.Details.(*dbmsOnlyPostgresqlInstanceDetails)
	if !ok {
		return nil, nil, fmt.Errorf(
			"error casting instance.Details " +
				"as *dbmsOnlyPostgresqlInstanceDetails",
		)
	}
	result, err := d.serversClient.Delete(
		ctx,
		instance.ResourceGroup,
		dt.ServerName,
	)
	if err != nil {
		return nil, nil, fmt.Errorf("error deleting postgresql server: %s", err)
	}
	if err := result.WaitForCompletion(ctx, d.serversClient.Client); err != nil {
		return nil, nil, fmt.Errorf("error deleting postgresql server: %s", err)
	}
	return dt, instance.SecureDetails, nil
}

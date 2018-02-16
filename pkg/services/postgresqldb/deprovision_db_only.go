package postgresqldb

import (
	"context"
	"fmt"

	"github.com/Azure/open-service-broker-azure/pkg/service"
)

func (d *dbOnlyManager) GetDeprovisioner(
	service.Plan,
) (service.Deprovisioner, error) {
	return service.NewDeprovisioner(
		service.NewDeprovisioningStep("deleteARMDeployment", d.deleteARMDeployment),
		service.NewDeprovisioningStep(
			"deletePostgreSQLDatabase",
			d.deletePostgreSQLDatabase,
		),
	)
}

func (d *dbOnlyManager) deleteARMDeployment(
	_ context.Context,
	instance service.Instance,
) (service.InstanceDetails, service.SecureInstanceDetails, error) {
	dt, ok := instance.Details.(*dbOnlyPostgresqlInstanceDetails)
	if !ok {
		return nil, nil, fmt.Errorf(
			"error casting instance.Details as *dbOnlyPostgresqlInstanceDetails",
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

func (d *dbOnlyManager) deletePostgreSQLDatabase(
	ctx context.Context,
	instance service.Instance,
) (service.InstanceDetails, service.SecureInstanceDetails, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	pdt, ok := instance.Parent.Details.(*dbmsOnlyPostgresqlInstanceDetails)
	if !ok {
		return nil, nil, fmt.Errorf(
			"error casting instance.Parent.Details " +
				"as *dbmsOnlyPostgresqlInstanceDetails",
		)
	}
	dt, ok := instance.Details.(*dbOnlyPostgresqlInstanceDetails)
	if !ok {
		return nil, nil, fmt.Errorf(
			"error casting instance.Details " +
				"as *dbOnlyPostgresqlInstanceDetails",
		)
	}
	result, err := d.databasesClient.Delete(
		ctx,
		instance.Parent.ResourceGroup,
		pdt.ServerName,
		dt.DatabaseName,
	)
	if err != nil {
		return nil, nil, fmt.Errorf("error deleting postgresql database: %s", err)
	}
	if err := result.WaitForCompletion(ctx, d.databasesClient.Client); err != nil {
		return nil, nil, fmt.Errorf("error deleting postgresql database: %s", err)
	}
	return dt, instance.SecureDetails, nil
}

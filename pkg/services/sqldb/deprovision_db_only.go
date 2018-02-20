package sqldb

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
			"deleteMsSQLDatabase",
			d.deleteMsSQLDatabase,
		),
	)
}

func (d *dbOnlyManager) deleteARMDeployment(
	_ context.Context,
	instance service.Instance,
) (service.InstanceDetails, service.SecureInstanceDetails, error) {
	dt, ok := instance.Details.(*mssqlDBOnlyInstanceDetails)
	if !ok {
		return nil, nil, fmt.Errorf(
			"error casting instance.Details as *mssqlDBOnlyInstanceDetails",
		)
	}
	//Parent should be set by the framework, but return an error if it is not set.
	if instance.Parent == nil {
		return nil, nil, fmt.Errorf("parent instance not set")
	}
	err := d.armDeployer.Delete(
		dt.ARMDeploymentName,
		instance.Parent.ResourceGroup,
	)
	if err != nil {
		return nil, nil, fmt.Errorf("error deleting ARM deployment: %s", err)
	}
	return dt, instance.SecureDetails, nil
}

func (d *dbOnlyManager) deleteMsSQLDatabase(
	ctx context.Context,
	instance service.Instance,
) (service.InstanceDetails, service.SecureInstanceDetails, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	dt, ok := instance.Details.(*mssqlDBOnlyInstanceDetails)
	if !ok {
		return nil, nil, fmt.Errorf(
			"error casting instance.Details as *mssqlDBOnlyInstanceDetails",
		)
	}
	//Parent should be set by the framework, but return an error if it is not set.
	if instance.Parent == nil {
		return nil, nil, fmt.Errorf("parent instance not set")
	}
	pdt, ok := instance.Parent.Details.(*mssqlVMOnlyInstanceDetails)
	if !ok {
		return nil, nil, fmt.Errorf(
			"error casting instance.Parent.Details as " +
				"*mssqlVMOnlyInstanceDetails",
		)
	}

	if _, err := d.databasesClient.Delete(
		ctx,
		instance.Parent.ResourceGroup,
		pdt.ServerName,
		dt.DatabaseName,
	); err != nil {
		return nil, nil, fmt.Errorf("error deleting sql database: %s", err)
	}
	return dt, instance.SecureDetails, nil
}

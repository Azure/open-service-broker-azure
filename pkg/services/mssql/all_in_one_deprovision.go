package mssql

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
		service.NewDeprovisioningStep(
			"deleteMsSQLServer",
			a.deleteMsSQLServer,
		),
	)
}

func (a *allInOneManager) deleteARMDeployment(
	_ context.Context,
	instance service.Instance,
) (service.InstanceDetails, error) {
	dt := instance.Details.(*allInOneInstanceDetails)
	err := a.armDeployer.Delete(
		dt.ARMDeploymentName,
		instance.ProvisioningParameters.GetString("resourceGroup"),
	)
	if err != nil {
		return nil, fmt.Errorf("error deleting ARM deployment: %s", err)
	}
	return instance.Details, nil
}

func (a *allInOneManager) deleteMsSQLServer(
	ctx context.Context,
	instance service.Instance,
) (service.InstanceDetails, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	dt := instance.Details.(*allInOneInstanceDetails)
	result, err := a.serversClient.Delete(
		ctx,
		instance.ProvisioningParameters.GetString("resourceGroup"),
		dt.ServerName,
	)
	if err != nil {
		return nil, fmt.Errorf("error deleting sql server: %s", err)
	}
	if err := result.WaitForCompletion(ctx, a.serversClient.Client); err != nil {
		return nil, fmt.Errorf("error deleting sql server: %s", err)
	}
	return instance.Details, nil
}

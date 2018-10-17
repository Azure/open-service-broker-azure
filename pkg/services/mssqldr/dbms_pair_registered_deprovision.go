package mssqldr

import (
	"context"

	"github.com/Azure/open-service-broker-azure/pkg/service"
)

func (d *dbmsPairRegisteredManager) GetDeprovisioner(
	service.Plan,
) (service.Deprovisioner, error) {
	return service.NewDeprovisioner(
		service.NewDeprovisioningStep("unregisterDBMSPair", d.unregisterDBMSPair),
	)
}

func (d *dbmsPairRegisteredManager) unregisterDBMSPair(
	_ context.Context,
	instance service.Instance,
) (service.InstanceDetails, error) {
	// do nothing, just for the framework to get the first step as it is required
	return instance.Details, nil
}

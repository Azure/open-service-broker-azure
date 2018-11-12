package mssqldr

import (
	"context"

	"github.com/Azure/open-service-broker-azure/pkg/service"
)

func (d *databasePairRegisteredManager) GetProvisioner(
	service.Plan,
) (service.Provisioner, error) {
	return service.NewProvisioner(
		service.NewProvisioningStep("preProvision", d.preProvision),
		service.NewProvisioningStep("validatePriDatabase", d.validatePriDatabase),
		service.NewProvisioningStep("validateSecDatabase", d.validateSecDatabase),
		service.NewProvisioningStep("validateFailoverGroup", d.validateFailoverGroup),
	)
}

func (d *databasePairRegisteredManager) preProvision(
	_ context.Context,
	instance service.Instance,
) (service.InstanceDetails, error) {
	pp := instance.ProvisioningParameters
	return &databasePairInstanceDetails{
		DatabaseName:      pp.GetString("database"),
		FailoverGroupName: pp.GetString("failoverGroup"),
	}, nil
}

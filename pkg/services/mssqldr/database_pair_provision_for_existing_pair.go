package mssqldr

import (
	"context"

	"github.com/Azure/open-service-broker-azure/pkg/service"
	uuid "github.com/satori/go.uuid"
)

func (d *databasePairManagerForExistingPair) GetProvisioner(
	service.Plan,
) (service.Provisioner, error) {
	return service.NewProvisioner(
		service.NewProvisioningStep("preProvision", d.preProvision),
		service.NewProvisioningStep("validatePriDatabase", d.validatePriDatabase),
		service.NewProvisioningStep("validateSecDatabase", d.validateSecDatabase),
		service.NewProvisioningStep(
			"validateFailoverGroup",
			d.validateFailoverGroup,
		),
		service.NewProvisioningStep(
			"deployPriARMTemplateForExistingInstance",
			d.deployPriARMTemplateForExistingInstance,
		),
		service.NewProvisioningStep(
			"deployFailoverGroupARMTemplateForExistingInstance",
			d.deployFailoverGroupARMTemplateForExistingInstance,
		),
		service.NewProvisioningStep(
			"deploySecARMTemplateForExistingInstance",
			d.deploySecARMTemplateForExistingInstance,
		),
	)
}

func (d *databasePairManagerForExistingPair) preProvision(
	_ context.Context,
	_ service.Instance,
) (service.InstanceDetails, error) {
	return &databasePairInstanceDetails{
		PriARMDeploymentName:           uuid.NewV4().String(),
		SecARMDeploymentName:           uuid.NewV4().String(),
		FailoverGroupARMDeploymentName: uuid.NewV4().String(),
	}, nil
}

package mssqldr

import (
	"context"
	"fmt"
	"strings"

	"github.com/Azure/open-service-broker-azure/pkg/service"
	uuid "github.com/satori/go.uuid"
)

func (d *databasePairManagerForExistingPrimary) GetProvisioner(
	service.Plan,
) (service.Provisioner, error) {
	return service.NewProvisioner(
		service.NewProvisioningStep("preProvision", d.preProvision),
		service.NewProvisioningStep(
			"checkNameAvailability",
			d.checkNameAvailability,
		),
		service.NewProvisioningStep("validatePriDatabase", d.validatePriDatabase),
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

func (d *databasePairManagerForExistingPrimary) preProvision(
	_ context.Context,
	_ service.Instance,
) (service.InstanceDetails, error) {
	return &databasePairInstanceDetails{
		PriARMDeploymentName:           uuid.NewV4().String(),
		SecARMDeploymentName:           uuid.NewV4().String(),
		FailoverGroupARMDeploymentName: uuid.NewV4().String(),
	}, nil
}

func (d *databasePairManagerForExistingPrimary) checkNameAvailability(
	ctx context.Context,
	instance service.Instance,
) (service.InstanceDetails, error) {
	pp := instance.ProvisioningParameters
	ppp := instance.Parent.ProvisioningParameters
	pdt := instance.Parent.Details.(*dbmsPairInstanceDetails)
	if err := validateDatabase(
		ctx,
		&d.databasesClient,
		ppp.GetString("secondaryResourceGroup"),
		pdt.SecServerName,
		pp.GetString("database"),
	); err != nil {
		if !strings.Contains(err.Error(), "ResourceNotFound") {
			return nil, err
		}
	} else {
		return nil, fmt.Errorf("Secondary database with the name is already " +
			"existed")
	}
	if err := validateFailoverGroup(
		ctx,
		&d.failoverGroupsClient,
		ppp.GetString("primaryResourceGroup"),
		pdt.PriServerName,
		pdt.SecServerName,
		pp.GetString("database"),
		pp.GetString("failoverGroup"),
	); err != nil {
		if !strings.Contains(err.Error(), "ResourceNotFound") {
			return nil, err
		}
	} else {
		return nil, fmt.Errorf("Failover group with the name is already existed")
	}
	return instance.Details, nil
}

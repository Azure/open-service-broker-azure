package mssqldr

import (
	"context"
	"fmt"
	"strings"

	"github.com/Azure/open-service-broker-azure/pkg/service"
	uuid "github.com/satori/go.uuid"
)

func (d *databasePairManager) GetProvisioner(
	service.Plan,
) (service.Provisioner, error) {
	return service.NewProvisioner(
		service.NewProvisioningStep("preProvision", d.preProvision),
		service.NewProvisioningStep(
			"checkNameAvailability",
			d.checkNameAvailability,
		),
		service.NewProvisioningStep(
			"deployPriARMTemplate",
			d.deployPriARMTemplate,
		),
		service.NewProvisioningStep(
			"deployFailoverGroupARMTemplate",
			d.deployFailoverGroupARMTemplate,
		),
		// The secondary database must be created by the creation of failover group.
		// This deployment is for the update api to update the secondary database.
		service.NewProvisioningStep(
			"deploySecARMTemplateForExistingInstance",
			d.deploySecARMTemplateForExistingInstance,
		),
	)
}

func (d *databasePairManager) preProvision(
	_ context.Context,
	_ service.Instance,
) (service.InstanceDetails, error) {
	return &databasePairInstanceDetails{
		PriARMDeploymentName:           uuid.NewV4().String(),
		SecARMDeploymentName:           uuid.NewV4().String(),
		FailoverGroupARMDeploymentName: uuid.NewV4().String(),
	}, nil
}

func (d *databasePairManager) checkNameAvailability(
	ctx context.Context,
	instance service.Instance,
) (service.InstanceDetails, error) {
	pp := instance.ProvisioningParameters
	ppp := instance.Parent.ProvisioningParameters
	pdt := instance.Parent.Details.(*dbmsPairInstanceDetails)
	if err := validateDatabase(
		ctx,
		&d.databasesClient,
		ppp.GetString("primaryResourceGroup"),
		pdt.PriServerName,
		pp.GetString("database"),
	); err != nil {
		if !strings.Contains(err.Error(), "ResourceNotFound") {
			return nil, err
		}
	} else {
		return nil, fmt.Errorf("Primary database with the name is already existed")
	}
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
		pp.GetString("failoverGroup"),
		pp.GetString("database"),
	); err != nil {
		if !strings.Contains(err.Error(), "ResourceNotFound") {
			return nil, err
		}
	} else {
		return nil, fmt.Errorf("Failover group with the name is already existed")
	}
	return instance.Details, nil
}

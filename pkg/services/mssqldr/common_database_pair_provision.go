package mssqldr

import (
	"context"
	"fmt"

	"github.com/Azure/open-service-broker-azure/pkg/service"
)

func (d *commonDatabasePairManager) validatePriDatabase(
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
		return nil, err
	}
	return instance.Details, nil
}

func (d *commonDatabasePairManager) validateSecDatabase(
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
		return nil, err
	}
	return instance.Details, nil
}

func (d *commonDatabasePairManager) validateFailoverGroup(
	ctx context.Context,
	instance service.Instance,
) (service.InstanceDetails, error) {
	pp := instance.ProvisioningParameters
	ppp := instance.Parent.ProvisioningParameters
	pdt := instance.Parent.Details.(*dbmsPairInstanceDetails)
	if err := validateFailoverGroup(
		ctx,
		&d.failoverGroupsClient,
		ppp.GetString("primaryResourceGroup"),
		pdt.PriServerName,
		pdt.SecServerName,
		pp.GetString("database"),
		pp.GetString("failoverGroup"),
	); err != nil {
		return nil, err
	}
	return instance.Details, nil
}

func (d *commonDatabasePairManager) deployPriARMTemplate(
	_ context.Context,
	instance service.Instance,
) (service.InstanceDetails, error) {
	dt := instance.Details.(*databasePairInstanceDetails)
	pdt := instance.Parent.Details.(*dbmsPairInstanceDetails)
	pp := instance.ProvisioningParameters
	ppp := instance.Parent.ProvisioningParameters
	pd := instance.Plan.GetProperties().Extended["tierDetails"].(planDetails)
	tagsObj := pp.GetObject("tags")
	tags := make(map[string]string, len(tagsObj.Data))
	for k := range tagsObj.Data {
		tags[k] = tagsObj.GetString(k)
	}
	if err := deployDatabaseARMTemplate(
		&d.armDeployer,
		dt.PriARMDeploymentName,
		ppp.GetString("primaryResourceGroup"),
		ppp.GetString("primaryLocation"),
		pdt.PriServerName,
		pp.GetString("database"),
		*pp,
		pd,
		tags,
	); err != nil {
		return nil, fmt.Errorf("error deploying ARM template: %s", err)
	}
	dt.DatabaseName = pp.GetString("database")
	return dt, nil
}

func (d *commonDatabasePairManager) deployFailoverGroupARMTemplate(
	_ context.Context,
	instance service.Instance,
) (service.InstanceDetails, error) {
	pp := instance.ProvisioningParameters
	if err := deployFailoverGroupARMTemplate(
		&d.armDeployer,
		instance,
	); err != nil {
		return nil, fmt.Errorf("error deploying ARM template: %s", err)
	}
	dt := instance.Details.(*databasePairInstanceDetails)
	dt.FailoverGroupName = pp.GetString("failoverGroup")
	return dt, nil
}

func (d *commonDatabasePairManager) deploySecARMTemplateForExistingInstance(
	_ context.Context,
	instance service.Instance,
) (service.InstanceDetails, error) {
	dt := instance.Details.(*databasePairInstanceDetails)
	pdt := instance.Parent.Details.(*dbmsPairInstanceDetails)
	pp := instance.ProvisioningParameters
	ppp := instance.Parent.ProvisioningParameters
	tagsObj := pp.GetObject("tags")
	tags := make(map[string]string, len(tagsObj.Data))
	for k := range tagsObj.Data {
		tags[k] = tagsObj.GetString(k)
	}
	if err := deployDatabaseARMTemplateForExistingInstance(
		&d.armDeployer,
		dt.SecARMDeploymentName,
		ppp.GetString("secondaryResourceGroup"),
		ppp.GetString("secondaryLocation"),
		pdt.SecServerName,
		pp.GetString("database"),
		tags,
	); err != nil {
		return nil, fmt.Errorf("error deploying ARM template: %s", err)
	}
	dt.DatabaseName = pp.GetString("database")
	dt.FailoverGroupName = pp.GetString("failoverGroup")
	return dt, nil
}

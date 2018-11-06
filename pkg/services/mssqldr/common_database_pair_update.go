package mssqldr

import (
	"context"
	"fmt"

	"github.com/Azure/open-service-broker-azure/pkg/service"
)

func (d *commonDatabasePairManager) ValidateUpdatingParameters(
	instance service.Instance,
) error {
	td := instance.Plan.GetProperties().Extended["tierDetails"]
	details := td.(planDetails)
	return details.validateUpdateParameters(instance)
}

func (d *commonDatabasePairManager) GetUpdater(
	service.Plan,
) (service.Updater, error) {
	return service.NewUpdater(
		service.NewUpdatingStep("updatePriARMTemplate", d.updatePriARMTemplate),
		service.NewUpdatingStep("updateSecARMTemplate", d.updateSecARMTemplate),
	)
}

func (d *commonDatabasePairManager) updatePriARMTemplate(
	_ context.Context,
	instance service.Instance,
) (service.InstanceDetails, error) {
	dt := instance.Details.(*databasePairInstanceDetails)
	pdt := instance.Parent.Details.(*dbmsPairInstanceDetails)
	ppp := instance.Parent.ProvisioningParameters
	pd := instance.Plan.GetProperties().Extended["tierDetails"].(planDetails)
	tagsObj := instance.ProvisioningParameters.GetObject("tags")
	tags := make(map[string]string, len(tagsObj.Data))
	for k := range tagsObj.Data {
		tags[k] = tagsObj.GetString(k)
	}
	err := updateDatabaseARMTemplate(
		&d.armDeployer,
		dt.PriARMDeploymentName,
		ppp.GetString("primaryResourceGroup"),
		ppp.GetString("primaryLocation"),
		pdt.PriServerName,
		dt.DatabaseName,
		*instance.UpdatingParameters,
		pd,
		tags,
	)
	if err != nil {
		return nil, fmt.Errorf("error deploying ARM template: %s", err)
	}
	return instance.Details, nil
}

func (d *commonDatabasePairManager) updateSecARMTemplate(
	_ context.Context,
	instance service.Instance,
) (service.InstanceDetails, error) {
	dt := instance.Details.(*databasePairInstanceDetails)
	pdt := instance.Parent.Details.(*dbmsPairInstanceDetails)
	ppp := instance.Parent.ProvisioningParameters
	pd := instance.Plan.GetProperties().Extended["tierDetails"].(planDetails)
	tagsObj := instance.ProvisioningParameters.GetObject("tags")
	tags := make(map[string]string, len(tagsObj.Data))
	for k := range tagsObj.Data {
		tags[k] = tagsObj.GetString(k)
	}
	err := updateDatabaseARMTemplate(
		&d.armDeployer,
		dt.SecARMDeploymentName,
		ppp.GetString("secondaryResourceGroup"),
		ppp.GetString("secondaryLocation"),
		pdt.SecServerName,
		dt.DatabaseName,
		*instance.UpdatingParameters,
		pd,
		tags,
	)
	if err != nil {
		return nil, fmt.Errorf("error deploying ARM template: %s", err)
	}
	return instance.Details, nil
}

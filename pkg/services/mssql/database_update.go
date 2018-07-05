package mssql

import (
	"context"
	"fmt"

	"github.com/Azure/open-service-broker-azure/pkg/service"
)

func (d *databaseManager) ValidateUpdatingParameters(
	instance service.Instance,
) error {
	td := instance.Plan.GetProperties().Extended["tierDetails"]
	details := td.(planDetails)
	return details.validateUpdateParameters(instance)
}

func (d *databaseManager) GetUpdater(service.Plan) (service.Updater, error) {
	// There isn't a need to do any "pre-provision here. just the update step"
	return service.NewUpdater(
		service.NewUpdatingStep("updateARMTemplate", d.updateARMTemplate),
	)
}

func (d *databaseManager) updateARMTemplate(
	_ context.Context,
	instance service.Instance,
) (service.InstanceDetails, error) {
	dt := instance.Details.(*databaseInstanceDetails)
	pdt := instance.Parent.Details.(*dbmsInstanceDetails)
	pd := instance.Plan.GetProperties().Extended["tierDetails"].(planDetails)
	goTemplateParams, err := buildDatabaseGoTemplateParameters(
		dt.DatabaseName,
		*instance.UpdatingParameters,
		pd,
	)
	if err != nil {
		return nil, err
	}
	goTemplateParams["location"] =
		instance.Parent.ProvisioningParameters.GetString("location")
	goTemplateParams["serverName"] = pdt.ServerName
	tagsObj := instance.ProvisioningParameters.GetObject("tags")
	tags := make(map[string]string, len(tagsObj.Data))
	for k := range tagsObj.Data {
		tags[k] = tagsObj.GetString(k)
	}
	_, err = d.armDeployer.Update(
		dt.ARMDeploymentName,
		instance.Parent.ProvisioningParameters.GetString("resourceGroup"),
		instance.Parent.ProvisioningParameters.GetString("location"),
		databaseARMTemplateBytes,
		goTemplateParams,
		map[string]interface{}{}, // empty arm template params
		tags,
	)
	if err != nil {
		return nil, fmt.Errorf("error deploying ARM template: %s", err)
	}
	// This shouldn't change the instance details, so just return
	// what was there already
	return instance.Details, err
}

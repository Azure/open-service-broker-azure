package mssql

import (
	"context"
	"fmt"

	"github.com/Azure/open-service-broker-azure/pkg/service"
)

func (a *allInOneManager) ValidateUpdatingParameters(
	instance service.Instance,
) error {
	td := instance.Plan.GetProperties().Extended["tierDetails"]
	details := td.(planDetails)
	return details.validateUpdateParameters(instance)
}

func (a *allInOneManager) GetUpdater(service.Plan) (service.Updater, error) {
	// There isn't a need to do any "pre-provision here. just the update step"
	return service.NewUpdater(
		service.NewUpdatingStep("updateARMTemplate", a.updateARMTemplate),
	)
}

func (a *allInOneManager) updateARMTemplate(
	_ context.Context,
	instance service.Instance,
) (service.InstanceDetails, error) {
	dt := instance.Details.(*allInOneInstanceDetails)
	version := instance.Service.GetProperties().Extended["version"].(string)
	goTemplateParams := buildDBMSGoTemplateParameters(
		&dt.dbmsInstanceDetails,
		*instance.UpdatingParameters,
		version,
	)
	pd := instance.Plan.GetProperties().Extended["tierDetails"].(planDetails)
	dbParams, err := buildDatabaseGoTemplateParameters(
		dt.DatabaseName,
		*instance.UpdatingParameters,
		pd,
	)
	if err != nil {
		return nil, err
	}
	for key, value := range dbParams {
		goTemplateParams[key] = value
	}
	goTemplateParams["location"] =
		instance.ProvisioningParameters.GetString("location")
	tagsObj := instance.ProvisioningParameters.GetObject("tags")
	tags := make(map[string]string, len(tagsObj.Data))
	for k := range tagsObj.Data {
		tags[k] = tagsObj.GetString(k)
	}
	_, err = a.armDeployer.Update(
		dt.ARMDeploymentName,
		instance.ProvisioningParameters.GetString("resourceGroup"),
		instance.ProvisioningParameters.GetString("location"),
		allInOneARMTemplateBytes,
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

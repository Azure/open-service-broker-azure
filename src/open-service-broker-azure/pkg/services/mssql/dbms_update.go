package mssql

import (
	"context"
	"fmt"

	"open-service-broker-azure/pkg/service"
)

func (d *dbmsManager) ValidateUpdatingParameters(service.Instance) error {
	return nil
}

func (d *dbmsManager) GetUpdater(service.Plan) (service.Updater, error) {
	// There isn't a need to do any "pre-provision here. just the update step"
	return service.NewUpdater(
		service.NewUpdatingStep("updateARMTemplate", d.updateARMTemplate),
	)
}

func (d *dbmsManager) updateARMTemplate(
	_ context.Context,
	instance service.Instance,
) (service.InstanceDetails, service.SecureInstanceDetails, error) {

	dt := dbmsInstanceDetails{}
	if err := service.GetStructFromMap(instance.Details, &dt); err != nil {
		return nil, nil, err
	}
	sdt := secureDBMSInstanceDetails{}
	if err :=
		service.GetStructFromMap(instance.SecureDetails, &sdt); err != nil {
		return nil, nil, err
	}
	version := instance.Service.GetProperties().Extended["version"].(string)
	goTemplateParams, err := buildDBMSGoTemplateParameters(
		dt,
		sdt,
		*instance.UpdatingParameters,
		version,
	)
	if err != nil {
		return nil, nil, fmt.Errorf("error deploying ARM template: %s", err)
	}
	goTemplateParams["location"] =
		instance.ProvisioningParameters.GetString("location")
	tagsObj := instance.ProvisioningParameters.GetObject("tags")
	tags := make(map[string]string, len(tagsObj.Data))
	for k := range tagsObj.Data {
		tags[k] = tagsObj.GetString(k)
	}
	_, err = d.armDeployer.Update(
		dt.ARMDeploymentName,
		instance.ProvisioningParameters.GetString("resourceGroup"),
		instance.ProvisioningParameters.GetString("location"),
		allInOneARMTemplateBytes,
		goTemplateParams,
		map[string]interface{}{}, // empty arm template params
		tags,
	)
	if err != nil {
		return nil, nil, fmt.Errorf("error deploying ARM template: %s", err)
	}

	// This shouldn't change the instance details, so just return
	// what was there already
	return instance.Details, instance.SecureDetails, err
}

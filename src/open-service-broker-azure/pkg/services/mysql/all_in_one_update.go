package mysql

import (
	"context"
	"fmt"

	"open-service-broker-azure/pkg/service"
)

func (a *allInOneManager) ValidateUpdatingParameters(
	instance service.Instance,
) error {
	return validateStorageUpdate(
		*instance.ProvisioningParameters,
		*instance.UpdatingParameters,
	)
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
) (service.InstanceDetails, service.SecureInstanceDetails, error) {
	dt := allInOneInstanceDetails{}
	if err := service.GetStructFromMap(instance.Details, &dt); err != nil {
		return nil, nil, err
	}

	sdt := secureAllInOneInstanceDetails{}
	if err := service.GetStructFromMap(instance.SecureDetails, &sdt); err != nil {
		return nil, nil, err
	}

	version := instance.Service.GetProperties().Extended["version"].(string)
	goTemplateParameters, err := buildGoTemplateParameters(
		instance.Plan,
		version,
		dt.dbmsInstanceDetails,
		sdt.secureDBMSInstanceDetails,
		*instance.UpdatingParameters,
	)

	if err != nil {
		return nil, nil, fmt.Errorf("unable to build go template parameters: %s", err)
	}
	goTemplateParameters["databaseName"] = dt.DatabaseName
	tagsObj := instance.UpdatingParameters.GetObject("tags")
	tags := make(map[string]string, len(tagsObj.Data))
	for k := range tagsObj.Data {
		tags[k] = tagsObj.GetString(k)
	}
	_, err = a.armDeployer.Update(
		dt.ARMDeploymentName,
		instance.UpdatingParameters.GetString("resourceGroup"),
		instance.UpdatingParameters.GetString("location"),
		dbmsARMTemplateBytes,
		goTemplateParameters,
		map[string]interface{}{},
		tags,
	)
	if err != nil {
		return nil, nil, fmt.Errorf("error deploying ARM template: %s", err)
	}

	// This shouldn't change the instance details, so just return
	// what was there already
	return instance.Details, instance.SecureDetails, err
}

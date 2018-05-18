package mysql

import (
	"context"
	"fmt"

	"github.com/Azure/open-service-broker-azure/pkg/service"
)

func (a *allInOneManager) ValidateUpdatingParameters(
	instance service.Instance,
) error {
	pp := dbmsProvisioningParameters{}
	if err := service.GetStructFromMap(
		instance.ProvisioningParameters,
		&pp,
	); err != nil {
		return err
	}
	up := dbmsUpdatingParameters{}
	if err := service.GetStructFromMap(
		instance.UpdatingParameters,
		&up,
	); err != nil {
		return err
	}
	return validateStorageUpdate(pp, up)
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

	up := dbmsProvisioningParameters{}
	if err := service.GetStructFromMap(
		instance.UpdatingParameters,
		&up,
	); err != nil {
		return nil, nil, err
	}

	version := instance.Service.GetProperties().Extended["version"].(string)
	goTemplateParameters, err := buildGoTemplateParameters(
		instance.Plan,
		version,
		dt.dbmsInstanceDetails,
		sdt.secureDBMSInstanceDetails,
		up,
	)

	if err != nil {
		return nil, nil, fmt.Errorf("unable to build go template parameters: %s", err)
	}
	goTemplateParameters["databaseName"] = dt.DatabaseName

	_, err = a.armDeployer.Update(
		dt.ARMDeploymentName,
		instance.ResourceGroup,
		instance.Location,
		dbmsARMTemplateBytes,
		goTemplateParameters,
		map[string]interface{}{},
		instance.Tags,
	)
	if err != nil {
		return nil, nil, fmt.Errorf("error deploying ARM template: %s", err)
	}

	// This shouldn't change the instance details, so just return
	// what was there already
	return instance.Details, instance.SecureDetails, err
}

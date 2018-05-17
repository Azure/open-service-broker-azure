package postgresql

import (
	"context"
	"fmt"

	"github.com/Azure/open-service-broker-azure/pkg/service"
)

func (d *dbmsManager) ValidateUpdatingParameters(
	instance service.Instance,
) error {
	pp := dbmsProvisioningParameters{}
	if err := service.GetStructFromMap(
		instance.ProvisioningParameters.Data,
		&pp,
	); err != nil {
		return err
	}
	up := dbmsUpdatingParameters{}
	if err := service.GetStructFromMap(
		instance.UpdatingParameters.Data,
		&up,
	); err != nil {
		return err
	}
	return validateDBMSUpdateParameters(instance.Plan, pp, up)
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
	if err := service.GetStructFromMap(instance.SecureDetails, &sdt); err != nil {
		return nil, nil, err
	}

	up := dbmsProvisioningParameters{}
	if err := service.GetStructFromMap(
		instance.UpdatingParameters.Data,
		&up,
	); err != nil {
		return nil, nil, err
	}

	version := instance.Service.GetProperties().Extended["version"].(string)
	goTemplateParameters, err := buildGoTemplateParameters(
		instance.Plan,
		version,
		dt,
		sdt,
		up,
	)

	if err != nil {
		return nil, nil, fmt.Errorf("unable to build go template parameters: %s", err)
	}

	_, err = d.armDeployer.Update(
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

	return instance.Details, instance.SecureDetails, err
}

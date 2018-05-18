package postgresql

import (
	"context"
	"fmt"

	"github.com/Azure/open-service-broker-azure/pkg/service"
)

func (d *dbmsManager) ValidateUpdatingParameters(
	instance service.Instance,
) error {
	return validateStorageUpdate(instance)
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

	goTemplateParameters, err := buildGoTemplateParameters(
		instance.Plan,
		dt,
		sdt,
		instance.UpdatingParameters,
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

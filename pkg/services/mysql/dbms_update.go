package mysql

import (
	"github.com/Azure/open-service-broker-azure/pkg/service"
)

func (d *dbmsManager) ValidateUpdatingParameters(
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

func (d *dbmsManager) GetUpdater(service.Plan) (service.Updater, error) {
	return service.NewUpdater()
}

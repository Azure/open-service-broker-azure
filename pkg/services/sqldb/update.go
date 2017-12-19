package sqldb

import (
	"github.com/Azure/open-service-broker-azure/pkg/service"
)

func (a *allServiceManager) ValidateUpdatingParameters(
	updatingParameters service.UpdatingParameters,
) error {
	return nil
}

func (a *allServiceManager) GetUpdater(service.Plan) (service.Updater, error) {
	return service.NewUpdater()
}

func (v *vmServiceManager) ValidateUpdatingParameters(
	updatingParameters service.UpdatingParameters,
) error {
	return nil
}

func (v *vmServiceManager) GetUpdater(service.Plan) (service.Updater, error) {
	return service.NewUpdater()
}

func (d *dbServiceManager) ValidateUpdatingParameters(
	updatingParameters service.UpdatingParameters,
) error {
	return nil
}

func (d *dbServiceManager) GetUpdater(service.Plan) (service.Updater, error) {
	return service.NewUpdater()
}

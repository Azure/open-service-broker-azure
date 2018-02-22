package sqldb

import (
	"github.com/Azure/open-service-broker-azure/pkg/service"
)

func (d *dbOnlyManager) ValidateUpdatingParameters(service.Instance) error {
	return nil
}

func (d *dbOnlyManager) GetUpdater(service.Plan) (service.Updater, error) {
	return service.NewUpdater()
}

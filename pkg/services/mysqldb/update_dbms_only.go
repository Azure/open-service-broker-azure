package mysqldb

import (
	"github.com/Azure/open-service-broker-azure/pkg/service"
)

func (v *dbmsOnlyManager) ValidateUpdatingParameters(service.Instance) error {
	return nil
}

func (v *dbmsOnlyManager) GetUpdater(service.Plan) (service.Updater, error) {
	return service.NewUpdater()
}

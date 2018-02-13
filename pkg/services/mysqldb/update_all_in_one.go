package mysqldb

import (
	"github.com/Azure/open-service-broker-azure/pkg/service"
)

func (a *allInOneManager) ValidateUpdatingParameters(
	updatingParameters service.UpdatingParameters,
) error {
	return nil
}

func (a *allInOneManager) GetUpdater(service.Plan) (service.Updater, error) {
	return service.NewUpdater()
}

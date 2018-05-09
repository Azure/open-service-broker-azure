package mysql

import (
	"github.com/Azure/open-service-broker-azure/pkg/service"
)

func (d *databaseManager) GetUpdater(service.Plan) (service.Updater, error) {
	return service.NewUpdater()
}

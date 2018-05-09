package mysql

import (
	"github.com/Azure/open-service-broker-azure/pkg/service"
)

func (v *dbmsManager) GetUpdater(service.Plan) (service.Updater, error) {
	return service.NewUpdater()
}

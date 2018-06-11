package mssql

import (
	"github.com/Azure/open-service-broker-azure/pkg/service"
)

func (d *dbmsManager) ValidateUpdatingParameters(service.Instance) error {
	return nil
}

func (d *dbmsManager) GetUpdater(service.Plan) (service.Updater, error) {
	return service.NewUpdater()
}

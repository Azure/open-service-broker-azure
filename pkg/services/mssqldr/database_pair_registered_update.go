package mssqldr

import (
	"github.com/Azure/open-service-broker-azure/pkg/service"
)

func (d *databasePairRegisteredManager) ValidateUpdatingParameters(
	service.Instance,
) error {
	return nil
}

func (d *databasePairRegisteredManager) GetUpdater(
	service.Plan,
) (service.Updater, error) {
	return service.NewUpdater()
}

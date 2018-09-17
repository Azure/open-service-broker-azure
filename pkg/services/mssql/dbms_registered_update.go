package mssql

import (
	"context"

	"github.com/Azure/open-service-broker-azure/pkg/service"
)

func (d *dbmsRegisteredManager) ValidateUpdatingParameters(
	service.Instance,
) error {
	return nil
}

func (d *dbmsRegisteredManager) GetUpdater(
	service.Plan,
) (service.Updater, error) {
	return service.NewUpdater(
		service.NewUpdatingStep("updateAdministrator", d.updateAdministrator),
		service.NewUpdatingStep("testConnection", d.testConnection),
	)
}

func (d *dbmsRegisteredManager) updateAdministrator(
	_ context.Context,
	instance service.Instance,
) (service.InstanceDetails, error) {
	dt := instance.Details.(*dbmsInstanceDetails)
	dt.AdministratorLogin =
		instance.UpdatingParameters.GetString("administratorLogin")
	dt.AdministratorLoginPassword = service.SecureString(
		instance.UpdatingParameters.GetString("administratorLoginPassword"),
	)
	return dt, nil
}

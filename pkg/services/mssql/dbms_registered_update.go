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
	)
}

func (d *dbmsRegisteredManager) updateAdministrator(
	_ context.Context,
	instance service.Instance,
) (service.InstanceDetails, error) {
	dt := instance.Details.(*dbmsInstanceDetails)
	updatedAdministratorLogin :=
		instance.UpdatingParameters.GetString("administratorLogin")
	updatedAdministratorLoginPassword :=
		instance.UpdatingParameters.GetString("administratorLoginPassword")

	if err := validateServerAdmin(
		updatedAdministratorLogin,
		updatedAdministratorLoginPassword,
		dt.FullyQualifiedDomainName,
	); err != nil {
		return nil, err
	}

	dt.AdministratorLogin = updatedAdministratorLogin
	dt.AdministratorLoginPassword = service.SecureString(
		updatedAdministratorLoginPassword,
	)

	return dt, nil
}

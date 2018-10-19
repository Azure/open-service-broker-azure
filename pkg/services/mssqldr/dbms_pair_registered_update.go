package mssqldr

import (
	"context"

	"github.com/Azure/open-service-broker-azure/pkg/service"
)

func (d *dbmsPairRegisteredManager) ValidateUpdatingParameters(
	service.Instance,
) error {
	return nil
}

func (d *dbmsPairRegisteredManager) GetUpdater(
	service.Plan,
) (service.Updater, error) {
	return service.NewUpdater(
		service.NewUpdatingStep("updateAdministrators", d.updateAdministrators),
	)
}

func (d *dbmsPairRegisteredManager) updateAdministrators(
	_ context.Context,
	instance service.Instance,
) (service.InstanceDetails, error) {
	dt := instance.Details.(*dbmsPairInstanceDetails)
	up := instance.UpdatingParameters
	updatedPriAdministratorLogin :=
		up.GetString("primaryAdministratorLogin")
	updatedPriAdministratorLoginPassword :=
		up.GetString("primaryAdministratorLoginPassword")
	if err := validateServerAdmin(
		updatedPriAdministratorLogin,
		updatedPriAdministratorLoginPassword,
		dt.PriFullyQualifiedDomainName,
	); err != nil {
		return nil, err
	}
	updatedSecAdministratorLogin :=
		up.GetString("secondaryAdministratorLogin")
	updatedSecAdministratorLoginPassword :=
		up.GetString("secondaryAdministratorLoginPassword")
	if err := validateServerAdmin(
		updatedSecAdministratorLogin,
		updatedSecAdministratorLoginPassword,
		dt.SecFullyQualifiedDomainName,
	); err != nil {
		return nil, err
	}

	dt.PriAdministratorLogin = updatedPriAdministratorLogin
	dt.PriAdministratorLoginPassword = service.SecureString(
		updatedPriAdministratorLoginPassword,
	)
	dt.SecAdministratorLogin = updatedSecAdministratorLogin
	dt.SecAdministratorLoginPassword = service.SecureString(
		updatedSecAdministratorLoginPassword,
	)
	return dt, nil
}

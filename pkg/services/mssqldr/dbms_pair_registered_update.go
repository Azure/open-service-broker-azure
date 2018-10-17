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
		service.NewUpdatingStep("testPriConnection", d.testPriConnection),
		service.NewUpdatingStep("testSecConnection", d.testSecConnection),
	)
}

func (d *dbmsPairRegisteredManager) updateAdministrators(
	_ context.Context,
	instance service.Instance,
) (service.InstanceDetails, error) {
	dt := instance.Details.(*dbmsPairInstanceDetails)
	up := instance.UpdatingParameters
	dt.PriAdministratorLogin = up.GetString("primaryAdministratorLogin")
	dt.PriAdministratorLoginPassword = service.SecureString(
		up.GetString("primaryAdministratorLoginPassword"), // nolint: lll
	)
	dt.SecAdministratorLogin = up.GetString("secondaryAdministratorLogin")
	dt.SecAdministratorLoginPassword = service.SecureString(
		up.GetString("secondaryAdministratorLoginPassword"), // nolint: lll
	)
	return dt, nil
}

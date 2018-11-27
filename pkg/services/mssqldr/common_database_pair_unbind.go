package mssqldr

import (
	"github.com/Azure/open-service-broker-azure/pkg/service"
)

func (d *commonDatabasePairManager) Unbind(
	instance service.Instance,
	binding service.Binding,
) error {
	dt := instance.Details.(*databasePairInstanceDetails)
	pdt := instance.Parent.Details.(*dbmsPairInstanceDetails)
	bd := binding.Details.(*bindingDetails)
	// TODO: detect who is the primary role now
	// assume the roles are not changed
	return unbind(
		pdt.PriAdministratorLogin,
		string(pdt.PriAdministratorLoginPassword),
		pdt.PriFullyQualifiedDomainName,
		dt.DatabaseName,
		bd,
	)
}

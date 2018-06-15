package mssql

import (
	"github.com/Azure/open-service-broker-azure/pkg/service"
)

func (d *databaseManager) Unbind(
	instance service.Instance,
	binding service.Binding,
) error {
	dt := instance.Details.(*databaseInstanceDetails)
	pdt := instance.Parent.Details.(*dbmsInstanceDetails)
	bd := binding.Details.(*bindingDetails)
	return unbind(
		pdt.AdministratorLogin,
		string(pdt.AdministratorLoginPassword),
		pdt.FullyQualifiedDomainName,
		dt.DatabaseName,
		bd,
	)
}

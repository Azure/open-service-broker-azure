package mssql

import (
	"fmt"

	"github.com/Azure/open-service-broker-azure/pkg/service"
)

func (d *databaseManager) Unbind(
	instance service.Instance,
	binding service.Binding,
) error {
	dt, ok := instance.Details.(*databaseInstanceDetails)
	if !ok {
		return fmt.Errorf(
			"error casting instance.Details as *mssql.databaseInstanceDetails",
		)
	}
	pdt, ok := instance.Parent.Details.(*dbmsInstanceDetails)
	if !ok {
		return fmt.Errorf(
			"error casting instance.Parent.Details as *mssql.dbmsInstanceDetails",
		)
	}
	spdt, ok := instance.Parent.SecureDetails.(*secureDBMSInstanceDetails)
	if !ok {
		return fmt.Errorf(
			"error casting instance.Parent.SecureDetails as" +
				"*mssql.secureDBMSInstanceDetails",
		)
	}
	bd, ok := binding.Details.(*bindingDetails)
	if !ok {
		return fmt.Errorf(
			"error casting binding.Details as *mssql.bindingDetails",
		)
	}
	return unbind(
		pdt.AdministratorLogin,
		spdt.AdministratorLoginPassword,
		pdt.FullyQualifiedDomainName,
		dt.DatabaseName,
		bd,
	)
}

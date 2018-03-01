package mysql

import (
	"fmt"

	"github.com/Azure/open-service-broker-azure/pkg/service"
)

func (d *databaseManager) Unbind(
	instance service.Instance,
	binding service.Binding,
) error {
	pdt, ok := instance.Parent.Details.(*dbmsInstanceDetails)
	if !ok {
		return fmt.Errorf(
			"error casting instance.Parent.Details as *mysql.dbmsInstanceDetails",
		)
	}
	spdt, ok := instance.Parent.SecureDetails.(*secureDBMSInstanceDetails)
	if !ok {
		return fmt.Errorf(
			"error casting instance.Parent.SecureDetails as " +
				"*mysql.secureDBMSInstanceDetails",
		)
	}
	dt, ok := instance.Details.(*databaseInstanceDetails)
	if !ok {
		return fmt.Errorf(
			"error casting instance.Details as *mysql.databaseInstanceDetails",
		)
	}
	bc, ok := binding.Details.(*bindingDetails)
	if !ok {
		return fmt.Errorf(
			"error casting binding.Details as *mysql.bindingDetails",
		)
	}

	return unbind(
		pdt.EnforceSSL,
		d.sqlDatabaseDNSSuffix,
		pdt.ServerName,
		spdt.AdministratorLoginPassword,
		pdt.FullyQualifiedDomainName,
		dt.DatabaseName,
		bc,
	)
}

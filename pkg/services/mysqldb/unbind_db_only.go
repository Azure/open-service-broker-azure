package mysqldb

import (
	"fmt"

	"github.com/Azure/open-service-broker-azure/pkg/service"
)

func (d *dbOnlyManager) Unbind(
	instance service.Instance,
	binding service.Binding,
) error {
	pdt, ok := instance.Parent.Details.(*dbmsOnlyMysqlInstanceDetails)
	if !ok {
		return fmt.Errorf(
			"error casting instance.Parent.Details as *dbmsOnlyMysqlInstanceDetails",
		)
	}
	spdt, ok := instance.Parent.SecureDetails.(*dbmsOnlyMysqlSecureInstanceDetails)
	if !ok {
		return fmt.Errorf(
			"error casting instance.Parent.SecureDetails as " +
				"*dbmsOnlyMysqlSecureInstanceDetails",
		)
	}
	dt, ok := instance.Details.(*dbOnlyMysqlInstanceDetails)
	if !ok {
		return fmt.Errorf(
			"error casting instance.Details as *dbOnlyMysqlInstanceDetails",
		)
	}
	bc, ok := binding.Details.(*mysqlBindingDetails)
	if !ok {
		return fmt.Errorf(
			"error casting binding.Details as *mysqlBindingDetails",
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

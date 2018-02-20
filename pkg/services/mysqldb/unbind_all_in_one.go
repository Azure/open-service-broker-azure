package mysqldb

import (
	"fmt"

	"github.com/Azure/open-service-broker-azure/pkg/service"
)

func (a *allInOneManager) Unbind(
	instance service.Instance,
	binding service.Binding,
) error {
	dt, ok := instance.Details.(*allInOneMysqlInstanceDetails)
	if !ok {
		return fmt.Errorf(
			"error casting instance.Details as *allInOneMysqlInstanceDetails",
		)
	}
	sdt, ok := instance.SecureDetails.(*allInOneMysqlSecureInstanceDetails)
	if !ok {
		return fmt.Errorf(
			"error casting instance.SecureDetails as " +
				"*allInOneMysqlSecureInstanceDetails",
		)
	}
	bc, ok := binding.Details.(*mysqlBindingDetails)
	if !ok {
		return fmt.Errorf(
			"error casting binding.Details as *mysqlBindingDetails",
		)
	}

	return unbind(
		dt.EnforceSSL,
		a.sqlDatabaseDNSSuffix,
		dt.ServerName,
		sdt.AdministratorLoginPassword,
		dt.FullyQualifiedDomainName,
		dt.DatabaseName,
		bc,
	)
}

package mysql

import (
	"fmt"

	"github.com/Azure/open-service-broker-azure/pkg/service"
)

func (a *allInOneManager) Unbind(
	instance service.Instance,
	binding service.Binding,
) error {
	dt, ok := instance.Details.(*allInOneInstanceDetails)
	if !ok {
		return fmt.Errorf(
			"error casting instance.Details as *mysql.allInOneInstanceDetails",
		)
	}
	sdt, ok := instance.SecureDetails.(*secureAllInOneInstanceDetails)
	if !ok {
		return fmt.Errorf(
			"error casting instance.SecureDetails as " +
				"*mysql.secureAllInOneInstanceDetails",
		)
	}
	bc, ok := binding.Details.(*bindingDetails)
	if !ok {
		return fmt.Errorf(
			"error casting binding.Details as *mysql.bindingDetails",
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

package postgresql

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
			"error casting instance.Parent.Details as " +
				"*postgresql.dbmsInstanceDetails",
		)
	}
	spdt, ok :=
		instance.Parent.SecureDetails.(*secureDBMSInstanceDetails)
	if !ok {
		return fmt.Errorf(
			"error casting instance.Parent.SecureDetails as " +
				"*postgresql.secureDBMSInstanceDetails",
		)
	}
	bd, ok := binding.Details.(*bindingDetails)
	if !ok {
		return fmt.Errorf(
			"error casting binding.Details as *postgresql.bindingDetails",
		)
	}

	return unbind(
		pdt.EnforceSSL,
		pdt.ServerName,
		spdt.AdministratorLoginPassword,
		pdt.FullyQualifiedDomainName,
		bd.LoginName,
	)
}

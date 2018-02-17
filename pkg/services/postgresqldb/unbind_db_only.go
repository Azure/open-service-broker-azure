package postgresqldb

import (
	"fmt"

	"github.com/Azure/open-service-broker-azure/pkg/service"
)

func (d *dbOnlyManager) Unbind(
	instance service.Instance,
	binding service.Binding,
) error {
	pdt, ok := instance.Parent.Details.(*dbmsOnlyPostgresqlInstanceDetails)
	if !ok {
		return fmt.Errorf(
			"error casting instance.Parent.Details " +
				"as *dbmsOnlyPostgresqlInstanceDetails",
		)
	}
	spdt, ok :=
		instance.Parent.SecureDetails.(*dbmsOnlyPostgresqlSecureInstanceDetails)
	if !ok {
		return fmt.Errorf(
			"error casting instance.Parent.SecureDetails " +
				"as *dbmsOnlyPostgresqlSecureInstanceDetails",
		)
	}
	bd, ok := binding.Details.(*postgresqlBindingDetails)
	if !ok {
		return fmt.Errorf(
			"error casting binding.Details as *postgresqlBindingDetails",
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

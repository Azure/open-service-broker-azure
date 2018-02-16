package postgresqldb

import (
	"fmt"

	"github.com/Azure/open-service-broker-azure/pkg/service"
)

func (d *dbOnlyManager) Unbind(
	instance service.Instance,
	bindingDetails service.BindingDetails,
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
	bc, ok := bindingDetails.(*postgresqlBindingDetails)
	if !ok {
		return fmt.Errorf(
			"error casting bindingDetails as *postgresqlBindingDetails",
		)
	}

	return unbind(
		pdt.EnforceSSL,
		pdt.ServerName,
		spdt.AdministratorLoginPassword,
		pdt.FullyQualifiedDomainName,
		bc,
	)
}

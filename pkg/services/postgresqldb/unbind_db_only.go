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
	bc, ok := bindingDetails.(*postgresqlBindingDetails)
	if !ok {
		return fmt.Errorf(
			"error casting bindingDetails as *postgresqlBindingDetails",
		)
	}

	return unbind(
		pdt.EnforceSSL,
		pdt.ServerName,
		pdt.AdministratorLoginPassword,
		pdt.FullyQualifiedDomainName,
		bc,
	)
}

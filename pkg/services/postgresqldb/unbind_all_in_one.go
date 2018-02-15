package postgresqldb

import (
	"fmt"

	"github.com/Azure/open-service-broker-azure/pkg/service"
)

func (a *allInOneManager) Unbind(
	instance service.Instance,
	bindingDetails service.BindingDetails,
) error {
	dt, ok := instance.Details.(*allInOnePostgresqlInstanceDetails)
	if !ok {
		return fmt.Errorf(
			"error casting instance.Details " +
				"as *allInOnePostgresqlInstanceDetails",
		)
	}
	sdt, ok := instance.SecureDetails.(*allInOnePostgresqlSecureInstanceDetails)
	if !ok {
		return fmt.Errorf(
			"error casting instance.SecureDetails " +
				"as *allInOnePostgresqlSecureInstanceDetails",
		)
	}
	bc, ok := bindingDetails.(*postgresqlBindingDetails)
	if !ok {
		return fmt.Errorf(
			"error casting bindingDetails as *postgresqlBindingDetails",
		)
	}

	return unbind(
		dt.EnforceSSL,
		dt.ServerName,
		sdt.AdministratorLoginPassword,
		dt.FullyQualifiedDomainName,
		bc,
	)
}

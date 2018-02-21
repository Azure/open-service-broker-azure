package sqldb

import (
	"fmt"

	"github.com/Azure/open-service-broker-azure/pkg/service"
)

func (a *allInOneManager) Unbind(
	instance service.Instance,
	binding service.Binding,
) error {
	dt, ok := instance.Details.(*mssqlAllInOneInstanceDetails)
	if !ok {
		return fmt.Errorf(
			"error casting instance.Details as *mssqlAllInOneInstanceDetails",
		)
	}
	sdt, ok := instance.SecureDetails.(*mssqlAllInOneSecureInstanceDetails)
	if !ok {
		return fmt.Errorf(
			"error casting instance.SecureDetails as " +
				"*mssqlAllInOneSecureInstanceDetails",
		)
	}
	bd, ok := binding.Details.(*mssqlBindingDetails)
	if !ok {
		return fmt.Errorf(
			"error casting binding.Details as *mssqlBindingDetails",
		)
	}

	return unbind(
		dt.AdministratorLogin,
		sdt.AdministratorLoginPassword,
		dt.FullyQualifiedDomainName,
		dt.DatabaseName,
		bd,
	)
}

package sqldb

import (
	"fmt"

	"github.com/Azure/open-service-broker-azure/pkg/service"
)

func (d *dbOnlyManager) Unbind(
	instance service.Instance,
	binding service.Binding,
) error {
	dt, ok := instance.Details.(*mssqlDBOnlyInstanceDetails)
	if !ok {
		return fmt.Errorf(
			"error casting instance.Details as *mssqlDBOnlyInstanceDetails",
		)
	}
	bd, ok := binding.Details.(*mssqlBindingDetails)
	if !ok {
		return fmt.Errorf(
			"error casting binding.Details as *mssqlBindingDetails",
		)
	}
	//Parent should be set by the framework, but return an error if it is not set.
	if instance.Parent == nil {
		return fmt.Errorf("parent instance not set")
	}
	pdt, ok := instance.Parent.Details.(*mssqlVMOnlyInstanceDetails)
	if !ok {
		return fmt.Errorf(
			"error casting instance.Parent.Details as" +
				"*mssqlVMOnlyInstanceDetails",
		)
	}
	spdt, ok := instance.Parent.SecureDetails.(*mssqlVMOnlySecureInstanceDetails)
	if !ok {
		return fmt.Errorf(
			"error casting instance.Parent.SecureDetails as" +
				"*mssqlVMOnlySecureInstanceDetails",
		)
	}

	return unbind(
		pdt.AdministratorLogin,
		spdt.AdministratorLoginPassword,
		dt.FullyQualifiedDomainName,
		dt.DatabaseName,
		bd,
	)
}

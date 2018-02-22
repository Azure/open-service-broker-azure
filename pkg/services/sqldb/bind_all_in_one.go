package sqldb

import (
	"fmt"

	"github.com/Azure/open-service-broker-azure/pkg/service"
)

func (a *allInOneManager) ValidateBindingParameters(
	bindingParameters service.BindingParameters,
) error {
	// There are no parameters for binding to MSSQL, so there is nothing
	// to validate
	return nil
}
func (a *allInOneManager) Bind(
	instance service.Instance,
	_ service.BindingParameters,
	_ service.SecureBindingParameters,
) (service.BindingDetails, service.SecureBindingDetails, error) {
	dt, ok := instance.Details.(*mssqlAllInOneInstanceDetails)
	if !ok {
		return nil, nil, fmt.Errorf(
			"error casting instance.Details as *mssqlAllInOneInstanceDetails",
		)
	}
	sdt, ok := instance.SecureDetails.(*mssqlAllInOneSecureInstanceDetails)
	if !ok {
		return nil, nil, fmt.Errorf(
			"error casting instance.SecureDetails as " +
				"*mssqlAllInOneSecureInstanceDetails",
		)
	}

	return bind(
		dt.AdministratorLogin,
		sdt.AdministratorLoginPassword,
		dt.FullyQualifiedDomainName,
		dt.DatabaseName,
	)
}

func (a *allInOneManager) GetCredentials(
	instance service.Instance,
	binding service.Binding,
) (service.Credentials, error) {
	dt, ok := instance.Details.(*mssqlAllInOneInstanceDetails)
	if !ok {
		return nil, fmt.Errorf(
			"error casting instance.Details as *mssqlAllInOneInstanceDetails",
		)
	}
	bd, ok := binding.Details.(*mssqlBindingDetails)
	if !ok {
		return nil, fmt.Errorf(
			"error casting binding.Details as *mssqlBindingDetails",
		)
	}
	sbd, ok := binding.SecureDetails.(*mssqlSecureBindingDetails)
	if !ok {
		return nil, fmt.Errorf(
			"error casting binding.SecureDetails as *mssqlSecureBindingDetails",
		)
	}
	creds := createCredential(
		dt.FullyQualifiedDomainName,
		dt.DatabaseName,
		bd.LoginName,
		sbd.Password,
	)
	return creds, nil
}

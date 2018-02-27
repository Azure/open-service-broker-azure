package mssql

import (
	"fmt"

	"github.com/Azure/open-service-broker-azure/pkg/service"
)

func (a *allInOneManager) ValidateBindingParameters(
	service.BindingParameters,
	service.SecureBindingParameters,
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
	dt, ok := instance.Details.(*allInOneInstanceDetails)
	if !ok {
		return nil, nil, fmt.Errorf(
			"error casting instance.Details as *mssql.allInOneInstanceDetails",
		)
	}
	sdt, ok := instance.SecureDetails.(*secureAllInOneInstanceDetails)
	if !ok {
		return nil, nil, fmt.Errorf(
			"error casting instance.SecureDetails as " +
				"*mssql.secureAllInOneInstanceDetails",
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
	dt, ok := instance.Details.(*allInOneInstanceDetails)
	if !ok {
		return nil, fmt.Errorf(
			"error casting instance.Details as *mssql.allInOneInstanceDetails",
		)
	}
	bd, ok := binding.Details.(*bindingDetails)
	if !ok {
		return nil, fmt.Errorf(
			"error casting binding.Details as *mssql.bindingDetails",
		)
	}
	sbd, ok := binding.SecureDetails.(*secureBindingDetails)
	if !ok {
		return nil, fmt.Errorf(
			"error casting binding.SecureDetails as *mssql.secureBindingDetails",
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

package mssql

import (
	"fmt"

	"github.com/Azure/open-service-broker-azure/pkg/service"
)

func (d *databaseManager) ValidateBindingParameters(
	service.BindingParameters,
	service.SecureBindingParameters,
) error {
	// There are no parameters for binding to MSSQL, so there is nothing
	// to validate
	return nil
}

func (d *databaseManager) Bind(
	instance service.Instance,
	_ service.BindingParameters,
	_ service.SecureBindingParameters,
) (service.BindingDetails, service.SecureBindingDetails, error) {

	dt, ok := instance.Details.(*databaseInstanceDetails)
	if !ok {
		return nil, nil, fmt.Errorf(
			"error casting instance.Details as *mssql.databaseInstanceDetails",
		)
	}
	// Parent should be set by the framework, but return an error if it is not
	// set.
	if instance.Parent == nil {
		return nil, nil, fmt.Errorf("parent instance not set")
	}
	pdt, ok := instance.Parent.Details.(*dbmsInstanceDetails)
	if !ok {
		return nil, nil, fmt.Errorf(
			"error casting instance.Parent.Details as *mssql.dbmsInstanceDetails",
		)
	}
	spdt, ok := instance.Parent.SecureDetails.(*secureDBMSInstanceDetails)
	if !ok {
		return nil, nil, fmt.Errorf(
			"error casting instance.Parent.SecureDetails as " +
				"*mssql.secureDBMSInstanceDetails",
		)
	}
	return bind(
		pdt.AdministratorLogin,
		spdt.AdministratorLoginPassword,
		pdt.FullyQualifiedDomainName,
		dt.DatabaseName,
	)
}

func (d *databaseManager) GetCredentials(
	instance service.Instance,
	binding service.Binding,
) (service.Credentials, error) {
	dt, ok := instance.Details.(*databaseInstanceDetails)
	if !ok {
		return nil, fmt.Errorf(
			"error casting instance.Details as *mssql.databaseInstanceDetails",
		)
	}
	pdt, ok := instance.Parent.Details.(*dbmsInstanceDetails)
	if !ok {
		return nil, fmt.Errorf(
			"error casting instance.Parent.Details as *mssql.dbmsInstanceDetails",
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
		pdt.FullyQualifiedDomainName,
		dt.DatabaseName,
		bd.LoginName,
		sbd.Password,
	)
	return creds, nil
}

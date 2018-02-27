package postgresql

import (
	"fmt"

	"github.com/Azure/open-service-broker-azure/pkg/service"
)

func (d *databaseManager) ValidateBindingParameters(
	service.BindingParameters,
	service.SecureBindingParameters,
) error {
	// There are no parameters for binding to PostgreSQL, so there is nothing
	// to validate
	return nil
}

func (d *databaseManager) Bind(
	instance service.Instance,
	_ service.BindingParameters,
	_ service.SecureBindingParameters,
) (service.BindingDetails, service.SecureBindingDetails, error) {
	pdt, ok := instance.Parent.Details.(*dbmsInstanceDetails)
	if !ok {
		return nil, nil, fmt.Errorf(
			"error casting instance.Parent.Details as " +
				"*postgresql.dbmsInstanceDetails",
		)
	}
	spdt, ok :=
		instance.Parent.SecureDetails.(*secureDBMSInstanceDetails)
	if !ok {
		return nil, nil, fmt.Errorf(
			"error casting instance.Parent.SecureDetails as " +
				"*postgresql.secureDBMSInstanceDetails",
		)
	}

	dt, ok := instance.Details.(*databaseInstanceDetails)
	if !ok {
		return nil, nil, fmt.Errorf(
			"error casting instance.Details as *postgresql.databaseInstanceDetails",
		)
	}

	bd, spd, err := createBinding(
		pdt.EnforceSSL,
		pdt.ServerName,
		spdt.AdministratorLoginPassword,
		pdt.FullyQualifiedDomainName,
		dt.DatabaseName,
	)
	return bd, spd, err
}

func (d *databaseManager) GetCredentials(
	instance service.Instance,
	binding service.Binding,
) (service.Credentials, error) {
	pdt, ok := instance.Parent.Details.(*dbmsInstanceDetails)
	if !ok {
		return nil, fmt.Errorf(
			"error casting instance.Parent.Details as " +
				"*postgresql.dbmsInstanceDetails",
		)
	}

	dt, ok := instance.Details.(*databaseInstanceDetails)
	if !ok {
		return nil, fmt.Errorf(
			"error casting instance.Details as *postgresql.databaseInstanceDetails",
		)
	}
	bd, ok := binding.Details.(*bindingDetails)
	if !ok {
		return nil, fmt.Errorf(
			"error casting binding.Details as *postgresql.bindingDetails",
		)
	}
	sbd, ok := binding.SecureDetails.(*secureBindingDetails)
	if !ok {
		return nil, fmt.Errorf(
			"error casting binding.SecureDetails as *postgresql.secureBindingDetails",
		)
	}
	cred := createCredential(
		pdt.FullyQualifiedDomainName,
		pdt.EnforceSSL,
		pdt.ServerName,
		dt.DatabaseName,
		bd,
		sbd,
	)
	return cred, nil
}

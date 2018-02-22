package postgresqldb

import (
	"fmt"

	"github.com/Azure/open-service-broker-azure/pkg/service"
)

func (d *dbOnlyManager) ValidateBindingParameters(
	bindingParameters service.BindingParameters,
) error {
	// There are no parameters for binding to PostgreSQL, so there is nothing
	// to validate
	return nil
}

func (d *dbOnlyManager) Bind(
	instance service.Instance,
	_ service.BindingParameters,
	_ service.SecureBindingParameters,
) (service.BindingDetails, service.SecureBindingDetails, error) {
	pdt, ok := instance.Parent.Details.(*dbmsOnlyPostgresqlInstanceDetails)
	if !ok {
		return nil, nil, fmt.Errorf(
			"error casting instance.Parent.Details " +
				"as *dbmsOnlyPostgresqlInstanceDetails",
		)
	}
	spdt, ok :=
		instance.Parent.SecureDetails.(*dbmsOnlyPostgresqlSecureInstanceDetails)
	if !ok {
		return nil, nil, fmt.Errorf(
			"error casting instance.Parent.SecureDetails " +
				"as *dbmsOnlyPostgresqlSecureInstanceDetails",
		)
	}

	dt, ok := instance.Details.(*dbOnlyPostgresqlInstanceDetails)
	if !ok {
		return nil, nil, fmt.Errorf(
			"error casting instance.Details " +
				"as *dbOnlyPostgresqlInstanceDetails",
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

func (d *dbOnlyManager) GetCredentials(
	instance service.Instance,
	binding service.Binding,
) (service.Credentials, error) {
	pdt, ok := instance.Parent.Details.(*dbmsOnlyPostgresqlInstanceDetails)
	if !ok {
		return nil, fmt.Errorf(
			"error casting instance.Parent.Details " +
				"as *dbmsOnlyPostgresqlInstanceDetails",
		)
	}

	dt, ok := instance.Details.(*dbOnlyPostgresqlInstanceDetails)
	if !ok {
		return nil, fmt.Errorf(
			"error casting instance.Details " +
				"as *dbOnlyPostgresqlInstanceDetails",
		)
	}
	bd, ok := binding.Details.(*postgresqlBindingDetails)
	if !ok {
		return nil, fmt.Errorf(
			"error casting binding.Details as *postgresqlBindingDetails",
		)
	}
	sbd, ok := binding.SecureDetails.(*postgresqlSecureBindingDetails)
	if !ok {
		return nil, fmt.Errorf(
			"error casting binding.SecureDetails as *postgresqlSecureBindingDetails",
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

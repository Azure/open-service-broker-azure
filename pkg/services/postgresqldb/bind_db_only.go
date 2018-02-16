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
) (service.BindingDetails, error) {
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

	binding, err := createBinding(
		pdt.EnforceSSL,
		pdt.ServerName,
		pdt.AdministratorLoginPassword,
		pdt.FullyQualifiedDomainName,
		dt.DatabaseName,
	)
	return binding, err
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
	return &Credentials{
		Host:     pdt.FullyQualifiedDomainName,
		Port:     5432,
		Database: dt.DatabaseName,
		Username: fmt.Sprintf("%s@%s", bd.LoginName, pdt.ServerName),
		Password: bd.Password,
	}, nil
}

package postgresqldb

import (
	"fmt"

	"github.com/Azure/open-service-broker-azure/pkg/service"
)

func (a *allInOneManager) ValidateBindingParameters(
	bindingParameters service.BindingParameters,
) error {
	// There are no parameters for binding to PostgreSQL, so there is nothing
	// to validate
	return nil
}

func (a *allInOneManager) Bind(
	instance service.Instance,
	_ service.BindingParameters,
) (service.BindingDetails, service.SecureBindingDetails, error) {
	dt, ok := instance.Details.(*allInOnePostgresqlInstanceDetails)
	if !ok {
		return nil, nil, fmt.Errorf(
			"error casting instance.Details " +
				"as *allInOnePostgresqlInstanceDetails",
		)
	}
	sdt, ok := instance.SecureDetails.(*allInOnePostgresqlSecureInstanceDetails)
	if !ok {
		return nil, nil, fmt.Errorf(
			"error casting instance.SecureDetails " +
				"as *allInOnePostgresqlSecureInstanceDetails",
		)
	}
	bd, spd, err := createBinding(
		dt.EnforceSSL,
		dt.ServerName,
		sdt.AdministratorLoginPassword,
		dt.FullyQualifiedDomainName,
		dt.DatabaseName,
	)
	return bd, spd, err
}

func (a *allInOneManager) GetCredentials(
	instance service.Instance,
	binding service.Binding,
) (service.Credentials, error) {
	dt, ok := instance.Details.(*allInOnePostgresqlInstanceDetails)
	if !ok {
		return nil, fmt.Errorf(
			"error casting instance.Details " +
				"as *allInOnePostgresqlInstanceDetails",
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
	return &Credentials{
		Host:     dt.FullyQualifiedDomainName,
		Port:     5432,
		Database: dt.DatabaseName,
		Username: fmt.Sprintf("%s@%s", bd.LoginName, dt.ServerName),
		Password: sbd.Password,
	}, nil
}

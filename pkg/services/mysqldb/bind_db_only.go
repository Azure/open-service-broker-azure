package mysqldb

import (
	"fmt"

	"github.com/Azure/open-service-broker-azure/pkg/service"
)

func (d *dbOnlyManager) ValidateBindingParameters(
	bindingParameters service.BindingParameters,
) error {
	// There are no parameters for binding to MySQL, so there is nothing
	// to validate
	return nil
}

func (d *dbOnlyManager) Bind(
	instance service.Instance,
	_ service.BindingParameters,
	_ service.SecureBindingParameters,
) (service.BindingDetails, service.SecureBindingDetails, error) {
	pdt, ok := instance.Parent.Details.(*dbmsOnlyMysqlInstanceDetails)
	if !ok {
		return nil, nil, fmt.Errorf(
			"error casting instance.Parent.Details " +
				"as *dbmsOnlyMysqlInstanceDetails",
		)
	}
	spdt, ok :=
		instance.Parent.SecureDetails.(*dbmsOnlyMysqlSecureInstanceDetails)
	if !ok {
		return nil, nil, fmt.Errorf(
			"error casting instance.Parent.SecureDetails " +
				"as *dbmsOnlyMysqlSecureInstanceDetails",
		)
	}
	dt, ok := instance.Details.(*dbOnlyMysqlInstanceDetails)
	if !ok {
		return nil, nil, fmt.Errorf(
			"error casting instance.Details as *dbOnlyMysqlInstanceDetails",
		)
	}

	return createBinding(
		pdt.EnforceSSL,
		d.sqlDatabaseDNSSuffix,
		pdt.ServerName,
		spdt.AdministratorLoginPassword,
		pdt.FullyQualifiedDomainName,
		dt.DatabaseName,
	)
}

func (d *dbOnlyManager) GetCredentials(
	instance service.Instance,
	binding service.Binding,
) (service.Credentials, error) {
	pdt, ok := instance.Parent.Details.(*dbmsOnlyMysqlInstanceDetails)
	if !ok {
		return nil, fmt.Errorf(
			"error casting instance.Parent.Details " +
				"as *dbmsOnlyMysqlInstanceDetails",
		)
	}
	dt, ok := instance.Details.(*dbOnlyMysqlInstanceDetails)
	if !ok {
		return nil, fmt.Errorf(
			"error casting instance.Details as *dbOnlyMysqlInstanceDetails",
		)
	}
	bd, ok := binding.Details.(*mysqlBindingDetails)
	if !ok {
		return nil, fmt.Errorf(
			"error casting binding.Details as *mysqlBindingDetails",
		)
	}
	sbd, ok := binding.SecureDetails.(*mysqlSecureBindingDetails)
	if !ok {
		return nil, fmt.Errorf(
			"error casting binding.SecureDetails as *mysqlSecureBindingDetails",
		)
	}
	creds := createCredential(
		pdt.FullyQualifiedDomainName,
		pdt.EnforceSSL,
		pdt.ServerName,
		dt.DatabaseName,
		bd,
		sbd,
	)
	return creds, nil
}

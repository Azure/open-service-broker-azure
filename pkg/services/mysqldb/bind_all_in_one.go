package mysqldb

import (
	"fmt"

	"github.com/Azure/open-service-broker-azure/pkg/service"
)

func (a *allInOneManager) ValidateBindingParameters(
	bindingParameters service.BindingParameters,
) error {
	// There are no parameters for binding to MySQL, so there is nothing
	// to validate
	return nil
}

func (a *allInOneManager) Bind(
	instance service.Instance,
	_ service.BindingParameters,
	_ service.SecureBindingParameters,
) (service.BindingDetails, service.SecureBindingDetails, error) {
	dt, ok := instance.Details.(*allInOneMysqlInstanceDetails)
	if !ok {
		return nil, nil, fmt.Errorf(
			"error casting instance.Details as *allInOneMysqlInstanceDetails",
		)
	}
	sdt, ok := instance.SecureDetails.(*allInOneMysqlSecureInstanceDetails)
	if !ok {
		return nil, nil, fmt.Errorf(
			"error casting instance.SecureDetails as " +
				"*allInOneMysqlSecureInstanceDetails",
		)
	}

	return createBinding(
		dt.EnforceSSL,
		a.sqlDatabaseDNSSuffix,
		dt.ServerName,
		sdt.AdministratorLoginPassword,
		dt.FullyQualifiedDomainName,
		dt.DatabaseName,
	)
}

func (a *allInOneManager) GetCredentials(
	instance service.Instance,
	binding service.Binding,
) (service.Credentials, error) {
	dt, ok := instance.Details.(*allInOneMysqlInstanceDetails)
	if !ok {
		return nil, fmt.Errorf(
			"error casting instance.Details as *allInOneMysqlInstanceDetails",
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
		dt.FullyQualifiedDomainName,
		dt.EnforceSSL,
		dt.ServerName,
		dt.DatabaseName,
		bd,
		sbd,
	)
	return creds, nil
}

package postgresql

import (
	"fmt"

	"github.com/Azure/open-service-broker-azure/pkg/service"
)

func (a *allInOneManager) ValidateBindingParameters(
	service.BindingParameters,
	service.SecureBindingParameters,
) error {
	// There are no parameters for binding to PostgreSQL, so there is nothing
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
			"error casting instance.Details as *postgresql.allInOneInstanceDetails",
		)
	}
	sdt, ok := instance.SecureDetails.(*secureAllInOneInstanceDetails)
	if !ok {
		return nil, nil, fmt.Errorf(
			"error casting instance.SecureDetails as " +
				"*postgresql.secureAllInOneInstanceDetails",
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
	dt, ok := instance.Details.(*allInOneInstanceDetails)
	if !ok {
		return nil, fmt.Errorf(
			"error casting instance.Details as *postgresql.allInOneInstanceDetails",
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
		dt.FullyQualifiedDomainName,
		dt.EnforceSSL,
		dt.ServerName,
		dt.DatabaseName,
		bd,
		sbd,
	)
	return cred, nil
}

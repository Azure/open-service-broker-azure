package postgresql

import (
	"github.com/Azure/open-service-broker-azure/pkg/service"
)

func (a *allInOneManager) Bind(
	instance service.Instance,
	_ service.BindingParameters,
) (service.BindingDetails, error) {
	dt := instance.Details.(*allInOneInstanceDetails)
	bd, err := createBinding(
		isSSLRequired(*instance.ProvisioningParameters),
		dt.AdministratorLogin,
		dt.ServerName,
		string(dt.AdministratorLoginPassword),
		dt.FullyQualifiedDomainName,
		dt.DatabaseName,
	)
	return bd, err
}

func (a *allInOneManager) GetCredentials(
	instance service.Instance,
	binding service.Binding,
) (service.Credentials, error) {
	dt := instance.Details.(*allInOneInstanceDetails)
	bd := binding.Details.(*bindingDetails)
	cred := createCredential(
		dt.FullyQualifiedDomainName,
		isSSLRequired(*instance.ProvisioningParameters),
		dt.ServerName,
		dt.DatabaseName,
		bd,
	)
	return cred, nil
}

package mysql

import (
	"github.com/Azure/open-service-broker-azure/pkg/service"
)

func (a *allInOneManager) Bind(
	instance service.Instance,
	bp service.BindingParameters,
) (service.BindingDetails, error) {
	dt := instance.Details.(*allInOneInstanceDetails)
	return createBinding(
		bp,
		isSSLRequired(*instance.ProvisioningParameters),
		a.sqlDatabaseDNSSuffix,
		dt.ServerName,
		dt.AdministratorLogin,
		string(dt.AdministratorLoginPassword),
		dt.FullyQualifiedDomainName,
		dt.DatabaseName,
	)
}

func (a *allInOneManager) GetCredentials(
	instance service.Instance,
	binding service.Binding,
) (service.Credentials, error) {
	dt := instance.Details.(*allInOneInstanceDetails)
	bd := binding.Details.(*bindingDetails)
	creds := createCredential(
		dt.FullyQualifiedDomainName,
		isSSLRequired(*instance.ProvisioningParameters),
		dt.ServerName,
		dt.DatabaseName,
		bd,
	)
	return creds, nil
}

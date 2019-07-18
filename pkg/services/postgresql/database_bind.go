package postgresql

import (
	"github.com/Azure/open-service-broker-azure/pkg/service"
)

func (d *databaseManager) Bind(
	instance service.Instance,
	_ service.BindingParameters,
) (service.BindingDetails, error) {
	pdt := instance.Parent.Details.(*dbmsInstanceDetails)
	dt := instance.Details.(*databaseInstanceDetails)
	bd, err := createBinding(
		isSSLRequired(*instance.Parent.ProvisioningParameters),
		pdt.AdministratorLogin,
		pdt.ServerName,
		string(pdt.AdministratorLoginPassword),
		pdt.FullyQualifiedDomainName,
		dt.DatabaseName,
	)
	return bd, err
}

func (d *databaseManager) GetCredentials(
	instance service.Instance,
	binding service.Binding,
) (service.Credentials, error) {
	pdt := instance.Parent.Details.(*dbmsInstanceDetails)
	dt := instance.Details.(*databaseInstanceDetails)
	bd := binding.Details.(*bindingDetails)
	cred := createCredential(
		pdt.FullyQualifiedDomainName,
		isSSLRequired(*instance.Parent.ProvisioningParameters),
		pdt.ServerName,
		dt.DatabaseName,
		bd,
	)
	return cred, nil
}

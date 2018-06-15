package mssql

import (
	"github.com/Azure/open-service-broker-azure/pkg/service"
)

func (d *databaseManager) Bind(
	instance service.Instance,
	_ service.BindingParameters,
) (service.BindingDetails, error) {
	dt := instance.Details.(*databaseInstanceDetails)
	pdt := instance.Parent.Details.(*dbmsInstanceDetails)
	return bind(
		pdt.AdministratorLogin,
		string(pdt.AdministratorLoginPassword),
		pdt.FullyQualifiedDomainName,
		dt.DatabaseName,
	)
}

func (d *databaseManager) GetCredentials(
	instance service.Instance,
	binding service.Binding,
) (service.Credentials, error) {
	dt := instance.Details.(*databaseInstanceDetails)
	pdt := instance.Parent.Details.(*dbmsInstanceDetails)
	bd := binding.Details.(*bindingDetails)
	return createCredential(
		pdt.FullyQualifiedDomainName,
		dt.DatabaseName,
		bd.LoginName,
		string(bd.Password),
	), nil
}

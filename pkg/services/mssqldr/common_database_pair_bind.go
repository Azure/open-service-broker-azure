package mssqldr

import (
	"github.com/Azure/open-service-broker-azure/pkg/service"
)

func (d *commonDatabasePairManager) Bind(
	instance service.Instance,
	_ service.BindingParameters,
) (service.BindingDetails, error) {
	dt := instance.Details.(*databasePairInstanceDetails)
	pdt := instance.Parent.Details.(*dbmsPairInstanceDetails)
	// TODO: detect who is the primary role now
	// assume the roles are not changed
	return bind(
		pdt.PriAdministratorLogin,
		string(pdt.PriAdministratorLoginPassword),
		pdt.PriFullyQualifiedDomainName,
		dt.DatabaseName,
	)
}

func (d *commonDatabasePairManager) GetCredentials(
	instance service.Instance,
	binding service.Binding,
) (service.Credentials, error) {
	dt := instance.Details.(*databasePairInstanceDetails)
	pdt := instance.Parent.Details.(*dbmsPairInstanceDetails)
	bd := binding.Details.(*bindingDetails)
	return createCredential(
		pdt.PriFullyQualifiedDomainName,
		dt.FailoverGroupName,
		dt.DatabaseName,
		bd.Username,
		string(bd.Password),
	), nil
}

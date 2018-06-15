package postgresql

import (
	"github.com/Azure/open-service-broker-azure/pkg/service"
)

func (d *databaseManager) Unbind(
	instance service.Instance,
	binding service.Binding,
) error {
	pdt := instance.Parent.Details.(*dbmsInstanceDetails)
	bd := binding.Details.(*bindingDetails)
	return unbind(
		isSSLRequired(*instance.Parent.ProvisioningParameters),
		pdt.ServerName,
		string(pdt.AdministratorLoginPassword),
		pdt.FullyQualifiedDomainName,
		bd.LoginName,
	)
}

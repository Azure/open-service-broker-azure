package mysql

import (
	"github.com/Azure/open-service-broker-azure/pkg/service"
)

func (d *databaseManager) Unbind(
	instance service.Instance,
	binding service.Binding,
) error {
	pdt := instance.Parent.Details.(*dbmsInstanceDetails)
	dt := instance.Details.(*databaseInstanceDetails)
	bd := binding.Details.(*bindingDetails)
	return unbind(
		isSSLRequired(*instance.Parent.ProvisioningParameters),
		d.sqlDatabaseDNSSuffix,
		pdt.ServerName,
		pdt.AdministratorLogin,
		string(pdt.AdministratorLoginPassword),
		pdt.FullyQualifiedDomainName,
		dt.DatabaseName,
		bd,
	)
}

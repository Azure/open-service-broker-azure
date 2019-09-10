package mysql

import (
	"github.com/Azure/open-service-broker-azure/pkg/service"
)

func (a *allInOneManager) Unbind(
	instance service.Instance,
	binding service.Binding,
) error {
	dt := instance.Details.(*allInOneInstanceDetails)
	bd := binding.Details.(*bindingDetails)
	return unbind(
		isSSLRequired(*instance.ProvisioningParameters),
		a.sqlDatabaseDNSSuffix,
		dt.ServerName,
		dt.AdministratorLogin,
		string(dt.AdministratorLoginPassword),
		dt.FullyQualifiedDomainName,
		dt.DatabaseName,
		bd,
	)
}

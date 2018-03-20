package postgresql

import (
	"github.com/Azure/open-service-broker-azure/pkg/service"
)

func (d *databaseManager) Unbind(
	instance service.Instance,
	binding service.Binding,
) error {
	pdt := dbmsInstanceDetails{}
	if err :=
		service.GetStructFromMap(instance.Parent.Details, &pdt); err != nil {
		return err
	}
	spdt := secureDBMSInstanceDetails{}
	if err :=
		service.GetStructFromMap(instance.Parent.SecureDetails, &spdt); err != nil {
		return err
	}
	bd := bindingDetails{}
	if err := service.GetStructFromMap(binding.Details, &bd); err != nil {
		return err
	}

	return unbind(
		pdt.EnforceSSL,
		pdt.ServerName,
		spdt.AdministratorLoginPassword,
		pdt.FullyQualifiedDomainName,
		bd.LoginName,
	)
}

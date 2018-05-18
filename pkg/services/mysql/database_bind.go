package mysql

import (
	"github.com/Azure/open-service-broker-azure/pkg/service"
)

func (d *databaseManager) Bind(
	instance service.Instance,
	_ service.BindingParameters,
	_ service.SecureBindingParameters,
) (service.BindingDetails, service.SecureBindingDetails, error) {
	pdt := dbmsInstanceDetails{}
	if err :=
		service.GetStructFromMap(instance.Parent.Details, &pdt); err != nil {
		return nil, nil, err
	}
	spdt := secureDBMSInstanceDetails{}
	if err :=
		service.GetStructFromMap(instance.Parent.SecureDetails, &spdt); err != nil {
		return nil, nil, err
	}
	dt := databaseInstanceDetails{}
	if err := service.GetStructFromMap(instance.Details, &dt); err != nil {
		return nil, nil, err
	}
	ppp := dbmsProvisioningParameters{}
	if err := service.GetStructFromMap(
		instance.Parent.ProvisioningParameters,
		&ppp,
	); err != nil {
		return nil, nil, err
	}
	td := instance.Parent.Plan.GetProperties().Extended["tierDetails"].(tierDetails) // nolint: lll
	return createBinding(
		td.isSSLRequired(ppp),
		d.sqlDatabaseDNSSuffix,
		pdt.ServerName,
		spdt.AdministratorLoginPassword,
		pdt.FullyQualifiedDomainName,
		dt.DatabaseName,
	)
}

func (d *databaseManager) GetCredentials(
	instance service.Instance,
	binding service.Binding,
) (service.Credentials, error) {
	pdt := dbmsInstanceDetails{}
	if err :=
		service.GetStructFromMap(instance.Parent.Details, &pdt); err != nil {
		return nil, err
	}
	dt := databaseInstanceDetails{}
	if err := service.GetStructFromMap(instance.Details, &dt); err != nil {
		return nil, err
	}
	bd := bindingDetails{}
	if err := service.GetStructFromMap(binding.Details, &bd); err != nil {
		return nil, err
	}
	sbd := secureBindingDetails{}
	if err := service.GetStructFromMap(binding.SecureDetails, &sbd); err != nil {
		return nil, err
	}
	ppp := dbmsProvisioningParameters{}
	if err := service.GetStructFromMap(
		instance.Parent.ProvisioningParameters,
		&ppp,
	); err != nil {
		return nil, err
	}
	td := instance.Parent.Plan.GetProperties().Extended["tierDetails"].(tierDetails) // nolint: lll
	creds := createCredential(
		pdt.FullyQualifiedDomainName,
		td.isSSLRequired(ppp),
		pdt.ServerName,
		dt.DatabaseName,
		bd,
		sbd,
	)
	return creds, nil
}

package mysql

import (
	"open-service-broker-azure/pkg/service"
)

func (a *allInOneManager) Bind(
	instance service.Instance,
	_ service.BindingParameters,
) (service.BindingDetails, service.SecureBindingDetails, error) {
	dt := allInOneInstanceDetails{}
	if err := service.GetStructFromMap(instance.Details, &dt); err != nil {
		return nil, nil, err
	}
	sdt := secureAllInOneInstanceDetails{}
	if err := service.GetStructFromMap(instance.SecureDetails, &sdt); err != nil {
		return nil, nil, err
	}
	return createBinding(
		isSSLRequired(*instance.ProvisioningParameters),
		a.sqlDatabaseDNSSuffix,
		dt.ServerName,
		sdt.AdministratorLoginPassword,
		dt.FullyQualifiedDomainName,
		dt.DatabaseName,
	)
}

func (a *allInOneManager) GetCredentials(
	instance service.Instance,
	binding service.Binding,
) (service.Credentials, error) {
	dt := allInOneInstanceDetails{}
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
	creds := createCredential(
		dt.FullyQualifiedDomainName,
		isSSLRequired(*instance.ProvisioningParameters),
		dt.ServerName,
		dt.DatabaseName,
		bd,
		sbd,
	)
	return creds, nil
}

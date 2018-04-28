package mysql

import (
	"github.com/Azure/open-service-broker-azure/pkg/service"
)

func (a *allInOneManager) ValidateBindingParameters(
	service.BindingParameters,
	service.SecureBindingParameters,
) error {
	// There are no parameters for binding to MySQL, so there is nothing
	// to validate
	return nil
}

func (a *allInOneManager) Bind(
	instance service.Instance,
	_ service.BindingParameters,
	_ service.SecureBindingParameters,
) (service.BindingDetails, service.SecureBindingDetails, error) {
	dt := allInOneInstanceDetails{}
	if err := service.GetStructFromMap(instance.Details, &dt); err != nil {
		return nil, nil, err
	}
	sdt := secureAllInOneInstanceDetails{}
	if err := service.GetStructFromMap(instance.SecureDetails, &sdt); err != nil {
		return nil, nil, err
	}
	pp := allInOneProvisioningParameters{}
	if err :=
		service.GetStructFromMap(instance.ProvisioningParameters, &pp); err != nil {
		return nil, nil, err
	}
	schema :=
		instance.Plan.GetProperties().Extended["provisionSchema"].(planSchema)
	return createBinding(
		schema.isSSLRequired(pp.dbmsProvisioningParameters),
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
	pp := allInOneProvisioningParameters{}
	if err :=
		service.GetStructFromMap(instance.ProvisioningParameters, &pp); err != nil {
		return nil, err
	}
	schema :=
		instance.Plan.GetProperties().Extended["provisionSchema"].(planSchema)
	creds := createCredential(
		dt.FullyQualifiedDomainName,
		schema.isSSLRequired(pp.dbmsProvisioningParameters),
		dt.ServerName,
		dt.DatabaseName,
		bd,
		sbd,
	)
	return creds, nil
}

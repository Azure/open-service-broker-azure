package mysqldb

import "github.com/Azure/open-service-broker-azure/pkg/service"

type dbmsOnlyMysqlInstanceDetails struct {
	ARMDeploymentName        string `json:"armDeployment"`
	ServerName               string `json:"server"`
	FullyQualifiedDomainName string `json:"fullyQualifiedDomainName"`
	EnforceSSL               bool   `json:"enforceSSL"`
}

type dbmsOnlyMysqlSecureInstanceDetails struct {
	AdministratorLoginPassword string `json:"administratorLoginPassword"`
}

func (
	d *dbmsOnlyManager,
) GetEmptyProvisioningParameters() service.ProvisioningParameters {
	return &ServerProvisioningParameters{}
}

func (
	d *dbmsOnlyManager,
) GetEmptyUpdatingParameters() service.UpdatingParameters {
	return &UpdatingParameters{}
}

func (
	d *dbmsOnlyManager,
) GetEmptyInstanceDetails() service.InstanceDetails {
	return &dbmsOnlyMysqlInstanceDetails{}
}

func (
	d *dbmsOnlyManager,
) GetEmptySecureInstanceDetails() service.SecureInstanceDetails {
	return &dbmsOnlyMysqlSecureInstanceDetails{}
}

func (
	d *dbmsOnlyManager,
) GetEmptyBindingParameters() service.BindingParameters {
	return &BindingParameters{}
}

func (
	d *dbmsOnlyManager,
) GetEmptyBindingDetails() service.BindingDetails {
	return &mysqlBindingDetails{}
}

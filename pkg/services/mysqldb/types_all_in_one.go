package mysqldb

import "github.com/Azure/open-service-broker-azure/pkg/service"

type allInOneMysqlInstanceDetails struct {
	ARMDeploymentName        string `json:"armDeployment"`
	ServerName               string `json:"server"`
	DatabaseName             string `json:"database"`
	FullyQualifiedDomainName string `json:"fullyQualifiedDomainName"`
	EnforceSSL               bool   `json:"enforceSSL"`
}

type allInOneMysqlSecureInstanceDetails struct {
	AdministratorLoginPassword string `json:"administratorLoginPassword"`
}

func (
	a *allInOneManager,
) GetEmptyProvisioningParameters() service.ProvisioningParameters {
	return &ServerProvisioningParameters{}
}

func (
	a *allInOneManager,
) GetEmptySecureProvisioningParameters() service.SecureProvisioningParameters {
	return &SecureServerProvisioningParameters{}
}

func (
	a *allInOneManager,
) GetEmptyUpdatingParameters() service.UpdatingParameters {
	return &UpdatingParameters{}
}

func (
	a *allInOneManager,
) GetEmptyInstanceDetails() service.InstanceDetails {
	return &allInOneMysqlInstanceDetails{}
}

func (
	a *allInOneManager,
) GetEmptySecureInstanceDetails() service.SecureInstanceDetails {
	return &allInOneMysqlSecureInstanceDetails{}
}

func (
	a *allInOneManager,
) GetEmptyBindingParameters() service.BindingParameters {
	return &BindingParameters{}
}

func (a *allInOneManager) GetEmptyBindingDetails() service.BindingDetails {
	return &mysqlBindingDetails{}
}

func (
	a *allInOneManager,
) GetEmptySecureBindingDetails() service.SecureBindingDetails {
	return &mysqlSecureBindingDetails{}
}

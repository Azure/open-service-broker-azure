package mysqldb

import "github.com/Azure/open-service-broker-azure/pkg/service"

type allInOneMysqlInstanceDetails struct {
	ARMDeploymentName          string `json:"armDeployment"`
	ServerName                 string `json:"server"`
	AdministratorLoginPassword string `json:"administratorLoginPassword"`
	DatabaseName               string `json:"database"`
	FullyQualifiedDomainName   string `json:"fullyQualifiedDomainName"`
	EnforceSSL                 bool   `json:"enforceSSL"`
}

func (
	a *allInOneManager,
) GetEmptyProvisioningParameters() service.ProvisioningParameters {
	return &ServerProvisioningParameters{}
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
) GetEmptyBindingParameters() service.BindingParameters {
	return &BindingParameters{}
}

func (a *allInOneManager) GetEmptyBindingDetails() service.BindingDetails {
	return &mysqlBindingDetails{}
}

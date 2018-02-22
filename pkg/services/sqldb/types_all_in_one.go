package sqldb

import "github.com/Azure/open-service-broker-azure/pkg/service"

type mssqlAllInOneInstanceDetails struct {
	ARMDeploymentName string `json:"armDeployment"`
	serverInstanceDetails
	DatabaseName string `json:"database"`
}

type mssqlAllInOneSecureInstanceDetails struct {
	AdministratorLoginPassword string `json:"administratorLoginPassword"`
}

func (
	a *allInOneManager,
) GetEmptyProvisioningParameters() service.ProvisioningParameters {
	return &ServerProvisioningParams{}
}

func (
	a *allInOneManager,
) GetEmptySecureProvisioningParameters() service.SecureProvisioningParameters {
	return &SecureServerProvisioningParams{}
}

func (
	a *allInOneManager,
) GetEmptyInstanceDetails() service.InstanceDetails {
	return &mssqlAllInOneInstanceDetails{}
}

func (
	a *allInOneManager,
) GetEmptySecureInstanceDetails() service.SecureInstanceDetails {
	return &mssqlAllInOneSecureInstanceDetails{}
}

func (
	a *allInOneManager,
) GetEmptyBindingParameters() service.BindingParameters {
	return &BindingParameters{}
}

func (
	a *allInOneManager,
) GetEmptySecureBindingParameters() service.SecureBindingParameters {
	return &SecureBindingParameters{}
}

func (a *allInOneManager) GetEmptyBindingDetails() service.BindingDetails {
	return &mssqlBindingDetails{}
}

func (
	a *allInOneManager,
) GetEmptySecureBindingDetails() service.SecureBindingDetails {
	return &mssqlSecureBindingDetails{}
}

package sqldb

import "github.com/Azure/open-service-broker-azure/pkg/service"

type mssqlVMOnlyInstanceDetails struct {
	ARMDeploymentName string `json:"armDeployment"`
	serverInstanceDetails
}

type mssqlVMOnlySecureInstanceDetails struct {
	AdministratorLoginPassword string `json:"administratorLoginPassword"`
}

func (
	v *vmOnlyManager,
) GetEmptyProvisioningParameters() service.ProvisioningParameters {
	return &ServerProvisioningParams{}
}

func (
	v *vmOnlyManager,
) GetEmptyUpdatingParameters() service.UpdatingParameters {
	return &UpdatingParameters{}
}

func (
	v *vmOnlyManager,
) GetEmptyInstanceDetails() service.InstanceDetails {
	return &mssqlVMOnlyInstanceDetails{}
}

func (
	v *vmOnlyManager,
) GetEmptySecureInstanceDetails() service.SecureInstanceDetails {
	return &mssqlVMOnlySecureInstanceDetails{}
}

func (
	v *vmOnlyManager,
) GetEmptyBindingParameters() service.BindingParameters {
	return &BindingParameters{}
}

func (v *vmOnlyManager) GetEmptyBindingDetails() service.BindingDetails {
	return &mssqlBindingDetails{}
}

func (
	v *vmOnlyManager,
) GetEmptySecureBindingDetails() service.SecureBindingDetails {
	return &mssqlSecureBindingDetails{}
}

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
	d *dbmsOnlyManager,
) GetEmptyProvisioningParameters() service.ProvisioningParameters {
	return &ServerProvisioningParams{}
}

func (
	d *dbmsOnlyManager,
) GetEmptySecureProvisioningParameters() service.SecureProvisioningParameters {
	return &SecureServerProvisioningParams{}
}

func (
	d *dbmsOnlyManager,
) GetEmptyUpdatingParameters() service.UpdatingParameters {
	return &UpdatingParameters{}
}

func (
	d *dbmsOnlyManager,
) GetEmptyInstanceDetails() service.InstanceDetails {
	return &mssqlVMOnlyInstanceDetails{}
}

func (
	d *dbmsOnlyManager,
) GetEmptySecureInstanceDetails() service.SecureInstanceDetails {
	return &mssqlVMOnlySecureInstanceDetails{}
}

func (
	d *dbmsOnlyManager,
) GetEmptyBindingParameters() service.BindingParameters {
	return &BindingParameters{}
}

func (d *dbmsOnlyManager) GetEmptyBindingDetails() service.BindingDetails {
	return &mssqlBindingDetails{}
}

func (
	d *dbmsOnlyManager,
) GetEmptySecureBindingDetails() service.SecureBindingDetails {
	return &mssqlSecureBindingDetails{}
}

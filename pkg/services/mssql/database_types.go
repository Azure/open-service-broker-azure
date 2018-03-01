package mssql

import "github.com/Azure/open-service-broker-azure/pkg/service"

type databaseInstanceDetails struct {
	ARMDeploymentName string `json:"armDeployment"`
	DatabaseName      string `json:"database"`
}

func (
	d *databaseManager,
) GetEmptyProvisioningParameters() service.ProvisioningParameters {
	return nil
}

func (
	d *databaseManager,
) GetEmptySecureProvisioningParameters() service.SecureProvisioningParameters {
	return nil
}

func (
	d *databaseManager,
) GetEmptyInstanceDetails() service.InstanceDetails {
	return &databaseInstanceDetails{}
}

func (
	d *databaseManager,
) GetEmptySecureInstanceDetails() service.SecureInstanceDetails {
	return nil
}

func (
	d *databaseManager,
) GetEmptyBindingParameters() service.BindingParameters {
	return nil
}

func (
	d *databaseManager,
) GetEmptySecureBindingParameters() service.SecureBindingParameters {
	return nil
}

func (
	d *databaseManager,
) GetEmptyBindingDetails() service.BindingDetails {
	return &bindingDetails{}
}

func (
	d *databaseManager,
) GetEmptySecureBindingDetails() service.SecureBindingDetails {
	return &secureBindingDetails{}
}

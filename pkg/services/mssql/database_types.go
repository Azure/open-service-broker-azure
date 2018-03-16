package mssql

import "github.com/Azure/open-service-broker-azure/pkg/service"

// DatabaseProvisioningParams encapsulates non-sensitive database
// MS SQL-specific provisioning options
type DatabaseProvisioningParams struct {
	DisableTDE bool `json:"disableTransparentDataEncryption"`
}

type databaseInstanceDetails struct {
	ARMDeploymentName         string `json:"armDeployment"`
	DatabaseName              string `json:"database"`
	TransparentDataEncryption bool   `json:"transparentDataEncryption"`
}

func (
	d *databaseManager,
) GetEmptyProvisioningParameters() service.ProvisioningParameters {
	return &DatabaseProvisioningParams{}
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

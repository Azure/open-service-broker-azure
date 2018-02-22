package sqldb

import "github.com/Azure/open-service-broker-azure/pkg/service"

// DBProvisioningParams encapsulates non-sensitive MSSQL-specific provisioning
// options
type DBProvisioningParams struct {
}

// SecureDBProvisioningParams encapsulates sensitive MSSQL-specific provisioning
// options
type SecureDBProvisioningParams struct{}

type mssqlDBOnlyInstanceDetails struct {
	ARMDeploymentName        string `json:"armDeployment"`
	FullyQualifiedDomainName string `json:"fullyQualifiedDomainName"`
	DatabaseName             string `json:"database"`
}

type mssqlDBOnlySecureInstanceDetails struct{}

func (
	d *dbOnlyManager,
) GetEmptyProvisioningParameters() service.ProvisioningParameters {
	return &DBProvisioningParams{}
}

func (
	d *dbOnlyManager,
) GetEmptySecureProvisioningParameters() service.SecureProvisioningParameters {
	return &SecureDBProvisioningParams{}
}

func (
	d *dbOnlyManager,
) GetEmptyInstanceDetails() service.InstanceDetails {
	return &mssqlDBOnlyInstanceDetails{}
}

func (
	d *dbOnlyManager,
) GetEmptySecureInstanceDetails() service.SecureInstanceDetails {
	return &mssqlDBOnlySecureInstanceDetails{}
}

func (
	d *dbOnlyManager,
) GetEmptyBindingParameters() service.BindingParameters {
	return &BindingParameters{}
}

func (
	d *dbOnlyManager,
) GetEmptySecureBindingParameters() service.SecureBindingParameters {
	return &SecureBindingParameters{}
}

func (
	d *dbOnlyManager,
) GetEmptyBindingDetails() service.BindingDetails {
	return &mssqlBindingDetails{}
}

func (
	d *dbOnlyManager,
) GetEmptySecureBindingDetails() service.SecureBindingDetails {
	return &mssqlSecureBindingDetails{}
}

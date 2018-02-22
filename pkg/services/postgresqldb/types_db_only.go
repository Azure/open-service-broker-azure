package postgresqldb

import "github.com/Azure/open-service-broker-azure/pkg/service"

// DatabaseProvisioningParameters encapsulates non-sensitive PostgreSQL-specific
// database provisioning options
type DatabaseProvisioningParameters struct {
	Extensions []string `json:"extensions"`
}

// SecureDatabaseProvisioningParameters encapsulates sensitive
// PostgreSQL-specific database provisioning options
type SecureDatabaseProvisioningParameters struct{}

type dbOnlyPostgresqlInstanceDetails struct {
	ARMDeploymentName string `json:"armDeployment"`
	DatabaseName      string `json:"database"`
}

type dbOnlyPostgresqlSecureInstanceDetails struct{}

func (
	d *dbOnlyManager,
) GetEmptyProvisioningParameters() service.ProvisioningParameters {
	return &DatabaseProvisioningParameters{}
}

func (
	d *dbOnlyManager,
) GetEmptySecureProvisioningParameters() service.SecureProvisioningParameters {
	return &SecureDatabaseProvisioningParameters{}
}

func (
	d *dbOnlyManager,
) GetEmptyUpdatingParameters() service.UpdatingParameters {
	return &UpdatingParameters{}
}

func (
	d *dbOnlyManager,
) GetEmptyInstanceDetails() service.InstanceDetails {
	return &dbOnlyPostgresqlInstanceDetails{}
}

func (
	d *dbOnlyManager,
) GetEmptySecureInstanceDetails() service.SecureInstanceDetails {
	return &dbOnlyPostgresqlSecureInstanceDetails{}
}

func (d *dbOnlyManager) GetEmptyBindingParameters() service.BindingParameters {
	return &BindingParameters{}
}

func (
	d *dbOnlyManager,
) GetEmptySecureBindingParameters() service.SecureBindingParameters {
	return &SecureBindingParameters{}
}

func (d *dbOnlyManager) GetEmptyBindingDetails() service.BindingDetails {
	return &postgresqlBindingDetails{}
}

func (
	d *dbOnlyManager,
) GetEmptySecureBindingDetails() service.SecureBindingDetails {
	return &postgresqlSecureBindingDetails{}
}

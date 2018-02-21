package mysqldb

import "github.com/Azure/open-service-broker-azure/pkg/service"

// DatabaseProvisioningParameters encapsulates non-sensitive MySQL-specific
// database provisioning options
type DatabaseProvisioningParameters struct{}

// DatabaseSecureProvisioningParameters encapsulates sensitive MySQL-specific
// database provisioning options
type DatabaseSecureProvisioningParameters struct{}

type dbOnlyMysqlInstanceDetails struct {
	ARMDeploymentName string `json:"armDeployment"`
	DatabaseName      string `json:"database"`
}

type dbOnlyMysqlSecureInstanceDetails struct{}

func (
	d *dbOnlyManager,
) GetEmptyProvisioningParameters() service.ProvisioningParameters {
	return &DatabaseProvisioningParameters{}
}

func (
	d *dbOnlyManager,
) GetEmptySecureProvisioningParameters() service.SecureProvisioningParameters {
	return &DatabaseSecureProvisioningParameters{}
}

func (
	d *dbOnlyManager,
) GetEmptyUpdatingParameters() service.UpdatingParameters {
	return &UpdatingParameters{}
}

func (
	d *dbOnlyManager,
) GetEmptyInstanceDetails() service.InstanceDetails {
	return &dbOnlyMysqlInstanceDetails{}
}

func (
	d *dbOnlyManager,
) GetEmptySecureInstanceDetails() service.SecureInstanceDetails {
	return &dbOnlyMysqlSecureInstanceDetails{}
}

func (d *dbOnlyManager) GetEmptyBindingParameters() service.BindingParameters {
	return &BindingParameters{}
}

func (d *dbOnlyManager) GetEmptyBindingDetails() service.BindingDetails {
	return &mysqlBindingDetails{}
}

func (
	d *dbOnlyManager,
) GetEmptySecureBindingDetails() service.SecureBindingDetails {
	return &mysqlSecureBindingDetails{}
}

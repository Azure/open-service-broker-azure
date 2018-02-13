package mysqldb

import "github.com/Azure/open-service-broker-azure/pkg/service"

// DatabaseProvisioningParameters encapsulates MySQL-specific
// database provisioning options
type DatabaseProvisioningParameters struct{}

type dbOnlyMysqlInstanceDetails struct {
	ARMDeploymentName string `json:"armDeployment"`
	DatabaseName      string `json:"database"`
}

func (
	d *dbOnlyManager,
) GetEmptyProvisioningParameters() service.ProvisioningParameters {
	return &DatabaseProvisioningParameters{}
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
func (d *dbOnlyManager) GetEmptyBindingParameters() service.BindingParameters {
	return &BindingParameters{}
}

func (d *dbOnlyManager) GetEmptyBindingDetails() service.BindingDetails {
	return &mysqlBindingDetails{}
}

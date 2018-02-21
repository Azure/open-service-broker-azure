package sqldb

import "github.com/Azure/open-service-broker-azure/pkg/service"

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
) GetEmptyUpdatingParameters() service.UpdatingParameters {
	return &UpdatingParameters{}
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
) GetEmptyBindingDetails() service.BindingDetails {
	return &mssqlBindingDetails{}
}

func (
	d *dbOnlyManager,
) GetEmptySecureBindingDetails() service.SecureBindingDetails {
	return &mssqlSecureBindingDetails{}
}

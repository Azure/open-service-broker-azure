package postgresqldb

import "github.com/Azure/open-service-broker-azure/pkg/service"

// AllInOneProvisioningParameters encapsulates both dbms and database
// PostgreSQL-specific provisioning options
type AllInOneProvisioningParameters struct {
	ServerProvisioningParameters
	DatabaseProvisioningParameters
}

type allInOnePostgresqlInstanceDetails struct {
	dbmsOnlyPostgresqlInstanceDetails
	DatabaseName string `json:"database"`
}

type allInOnePostgresqlSecureInstanceDetails struct {
	dbmsOnlyPostgresqlSecureInstanceDetails
}

func (
	a *allInOneManager,
) GetEmptyProvisioningParameters() service.ProvisioningParameters {
	return &AllInOneProvisioningParameters{}
}

func (
	a *allInOneManager,
) GetEmptyUpdatingParameters() service.UpdatingParameters {
	return &UpdatingParameters{}
}

func (
	a *allInOneManager,
) GetEmptyInstanceDetails() service.InstanceDetails {
	return &allInOnePostgresqlInstanceDetails{}
}

func (
	a *allInOneManager,
) GetEmptySecureInstanceDetails() service.SecureInstanceDetails {
	return &allInOnePostgresqlSecureInstanceDetails{}
}

func (
	a *allInOneManager,
) GetEmptyBindingParameters() service.BindingParameters {
	return &BindingParameters{}
}

func (
	a *allInOneManager,
) GetEmptyBindingDetails() service.BindingDetails {
	return &postgresqlBindingDetails{}
}

func (
	a *allInOneManager,
) GetEmptySecureBindingDetails() service.SecureBindingDetails {
	return &postgresqlSecureBindingDetails{}
}

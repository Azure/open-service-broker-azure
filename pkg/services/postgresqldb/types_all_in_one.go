package postgresqldb

import "github.com/Azure/open-service-broker-azure/pkg/service"

// AllInOneProvisioningParameters encapsulates non-sensitive dbms AND database
// PostgreSQL-specific provisioning options
type AllInOneProvisioningParameters struct {
	ServerProvisioningParameters
	DatabaseProvisioningParameters
}

// SecureAllInOneProvisioningParameters encapsulates sensitive dbms AND database
// PostgreSQL-specific provisioning options
type SecureAllInOneProvisioningParameters struct{}

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
) GetEmptySecureProvisioningParameters() service.SecureProvisioningParameters {
	return &SecureAllInOneProvisioningParameters{}
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
) GetEmptySecureBindingParameters() service.SecureBindingParameters {
	return &SecureBindingParameters{}
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

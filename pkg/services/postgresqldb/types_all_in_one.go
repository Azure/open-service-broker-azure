package postgresqldb

import "github.com/Azure/open-service-broker-azure/pkg/service"

// AllInOneProvisioningParameters encapsulates both dbms and database
// PostgreSQL-specific provisioning options
type AllInOneProvisioningParameters struct {
	SSLEnforcement  string   `json:"sslEnforcement"`
	FirewallIPStart string   `json:"firewallStartIPAddress"`
	FirewallIPEnd   string   `json:"firewallEndIPAddress"`
	Extensions      []string `json:"extensions"`
}

type allInOnePostgresqlInstanceDetails struct {
	dbmsOnlyPostgresqlInstanceDetails
	DatabaseName string `json:"database"`
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
) GetEmptyBindingParameters() service.BindingParameters {
	return &BindingParameters{}
}

func (
	a *allInOneManager,
) GetEmptyBindingDetails() service.BindingDetails {
	return &postgresqlBindingDetails{}
}

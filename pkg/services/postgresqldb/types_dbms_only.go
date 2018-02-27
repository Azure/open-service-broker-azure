package postgresqldb

import "github.com/Azure/open-service-broker-azure/pkg/service"

// ServerProvisioningParameters encapsulates non-senstivie PostgreSQL-specific
// dbms provisioning options
type ServerProvisioningParameters struct {
	SSLEnforcement string         `json:"sslEnforcement"`
	FirewallRules  []FirewallRule `json:"firewallRules"`
}

// FirewallRule describes a firewall rule to be applied to an DBMS
type FirewallRule struct {
	Name    string `json:"name"`
	StartIP string `json:"startIPAddress"`
	EndIP   string `json:"endIPAddress"`
}

// SecureServerProvisioningParameters encapsulates senstivie PostgreSQL-specific
// dbms provisioning options
type SecureServerProvisioningParameters struct{}

type dbmsOnlyPostgresqlInstanceDetails struct {
	ARMDeploymentName        string `json:"armDeployment"`
	ServerName               string `json:"server"`
	FullyQualifiedDomainName string `json:"fullyQualifiedDomainName"`
	EnforceSSL               bool   `json:"enforceSSL"`
}

type dbmsOnlyPostgresqlSecureInstanceDetails struct {
	AdministratorLoginPassword string `json:"administratorLoginPassword"`
}

func (
	d *dbmsOnlyManager,
) GetEmptyProvisioningParameters() service.ProvisioningParameters {
	return &ServerProvisioningParameters{}
}

func (
	d *dbmsOnlyManager,
) GetEmptySecureProvisioningParameters() service.SecureProvisioningParameters {
	return &SecureServerProvisioningParameters{}
}

func (
	d *dbmsOnlyManager,
) GetEmptyInstanceDetails() service.InstanceDetails {
	return &dbmsOnlyPostgresqlInstanceDetails{}
}

func (
	d *dbmsOnlyManager,
) GetEmptySecureInstanceDetails() service.SecureInstanceDetails {
	return &dbmsOnlyPostgresqlSecureInstanceDetails{}
}

func (
	d *dbmsOnlyManager,
) GetEmptyBindingParameters() service.BindingParameters {
	return &BindingParameters{}
}

func (
	d *dbmsOnlyManager,
) GetEmptySecureBindingParameters() service.SecureBindingParameters {
	return &SecureBindingParameters{}
}

func (d *dbmsOnlyManager) GetEmptyBindingDetails() service.BindingDetails {
	return &postgresqlBindingDetails{}
}

func (
	d *dbmsOnlyManager,
) GetEmptySecureBindingDetails() service.SecureBindingDetails {
	return &postgresqlSecureBindingDetails{}
}

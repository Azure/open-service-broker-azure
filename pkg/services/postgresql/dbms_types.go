package postgresql

import "github.com/Azure/open-service-broker-azure/pkg/service"

// DBMSProvisioningParameters encapsulates non-senstivie PostgreSQL-specific
// dbms provisioning options
type DBMSProvisioningParameters struct {
	SSLEnforcement string         `json:"sslEnforcement"`
	FirewallRules  []FirewallRule `json:"firewallRules"`
}

// FirewallRule describes a firewall rule to be applied to an DBMS
type FirewallRule struct {
	Name    string `json:"name"`
	StartIP string `json:"startIPAddress"`
	EndIP   string `json:"endIPAddress"`
}

type dbmsInstanceDetails struct {
	ARMDeploymentName        string `json:"armDeployment"`
	ServerName               string `json:"server"`
	FullyQualifiedDomainName string `json:"fullyQualifiedDomainName"`
	EnforceSSL               bool   `json:"enforceSSL"`
}

type secureDBMSInstanceDetails struct {
	AdministratorLoginPassword string `json:"administratorLoginPassword"`
}

func (
	d *dbmsManager,
) GetEmptyProvisioningParameters() service.ProvisioningParameters {
	return &DBMSProvisioningParameters{}
}

func (
	d *dbmsManager,
) GetEmptySecureProvisioningParameters() service.SecureProvisioningParameters {
	return nil
}

func (
	d *dbmsManager,
) GetEmptyInstanceDetails() service.InstanceDetails {
	return &dbmsInstanceDetails{}
}

func (
	d *dbmsManager,
) GetEmptySecureInstanceDetails() service.SecureInstanceDetails {
	return &secureDBMSInstanceDetails{}
}

func (
	d *dbmsManager,
) GetEmptyBindingParameters() service.BindingParameters {
	return nil
}

func (
	d *dbmsManager,
) GetEmptySecureBindingParameters() service.SecureBindingParameters {
	return nil
}

func (d *dbmsManager) GetEmptyBindingDetails() service.BindingDetails {
	return nil
}

func (
	d *dbmsManager,
) GetEmptySecureBindingDetails() service.SecureBindingDetails {
	return &secureBindingDetails{}
}

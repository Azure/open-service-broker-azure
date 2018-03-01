package mssql

import "github.com/Azure/open-service-broker-azure/pkg/service"

// DBMSProvisioningParams encapsulates non-sensitive
// MSSQL-server specific provisioning options
type DBMSProvisioningParams struct {
	SSLEnforcement string         `json:"sslEnforcement"`
	FirewallRules  []FirewallRule `json:"firewallRules"`
}

// FirewallRule represents a firewall rule to be applied to the DBMS
type FirewallRule struct {
	Name    string `json:"name"`
	StartIP string `json:"startIPAddress"`
	EndIP   string `json:"endIPAddress"`
}

type dbmsInstanceDetails struct {
	ARMDeploymentName        string `json:"armDeployment"`
	FullyQualifiedDomainName string `json:"fullyQualifiedDomainName"`
	ServerName               string `json:"server"`
	AdministratorLogin       string `json:"administratorLogin"`
}

type secureDBMSInstanceDetails struct {
	AdministratorLoginPassword string `json:"administratorLoginPassword"`
}

func (
	d *dbmsManager,
) GetEmptyProvisioningParameters() service.ProvisioningParameters {
	return &DBMSProvisioningParams{}
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
	return &bindingDetails{}
}

func (
	d *dbmsManager,
) GetEmptySecureBindingDetails() service.SecureBindingDetails {
	return &secureBindingDetails{}
}

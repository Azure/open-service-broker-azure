package mysqldb

import "github.com/Azure/open-service-broker-azure/pkg/service"

// ServerProvisioningParameters encapsulates MySQL-specific
// server provisioning options
type ServerProvisioningParameters struct {
	SSLEnforcement  string `json:"sslEnforcement"`
	FirewallIPStart string `json:"firewallStartIPAddress"`
	FirewallIPEnd   string `json:"firewallEndIPAddress"`
}

// DatabaseProvisioningParameters encapsulates MySQL-specific
// database provisioning options
type DatabaseProvisioningParameters struct{}

type allInOneMysqlInstanceDetails struct {
	ARMDeploymentName          string `json:"armDeployment"`
	ServerName                 string `json:"server"`
	AdministratorLoginPassword string `json:"administratorLoginPassword"`
	DatabaseName               string `json:"database"`
	FullyQualifiedDomainName   string `json:"fullyQualifiedDomainName"`
	EnforceSSL                 bool   `json:"enforceSSL"`
}

type vmOnlyMysqlInstanceDetails struct {
	ARMDeploymentName          string `json:"armDeployment"`
	ServerName                 string `json:"server"`
	AdministratorLoginPassword string `json:"administratorLoginPassword"`
	FullyQualifiedDomainName   string `json:"fullyQualifiedDomainName"`
	EnforceSSL                 bool   `json:"enforceSSL"`
}

type dbOnlyMysqlInstanceDetails struct {
	ARMDeploymentName string `json:"armDeployment"`
	DatabaseName      string `json:"database"`
}

// UpdatingParameters encapsulates MySQL-specific updating options
type UpdatingParameters struct {
}

// BindingParameters encapsulates MySQL-specific binding options
type BindingParameters struct {
}

type mysqlBindingDetails struct {
	LoginName string `json:"loginName"`
	Password  string `json:"password"`
}

// Credentials encapsulates MySQL-specific coonection details and credentials.
type Credentials struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Database string `json:"database"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func (
	a *allInOneManager,
) GetEmptyProvisioningParameters() service.ProvisioningParameters {
	return &ServerProvisioningParameters{}
}

func (
	v *vmOnlyManager,
) GetEmptyProvisioningParameters() service.ProvisioningParameters {
	return &ServerProvisioningParameters{}
}

func (
	d *dbOnlyManager,
) GetEmptyProvisioningParameters() service.ProvisioningParameters {
	return &DatabaseProvisioningParameters{}
}

func (
	a *allInOneManager,
) GetEmptyUpdatingParameters() service.UpdatingParameters {
	return &UpdatingParameters{}
}

func (
	d *dbOnlyManager,
) GetEmptyUpdatingParameters() service.UpdatingParameters {
	return &UpdatingParameters{}
}

func (
	v *vmOnlyManager,
) GetEmptyUpdatingParameters() service.UpdatingParameters {
	return &UpdatingParameters{}
}

func (
	a *allInOneManager,
) GetEmptyInstanceDetails() service.InstanceDetails {
	return &allInOneMysqlInstanceDetails{}
}

func (
	v *vmOnlyManager,
) GetEmptyInstanceDetails() service.InstanceDetails {
	return &vmOnlyMysqlInstanceDetails{}
}

func (
	d *dbOnlyManager,
) GetEmptyInstanceDetails() service.InstanceDetails {
	return &dbOnlyMysqlInstanceDetails{}
}

func (
	a *allInOneManager,
) GetEmptyBindingParameters() service.BindingParameters {
	return &BindingParameters{}
}

func (v *vmOnlyManager) GetEmptyBindingParameters() service.BindingParameters {
	return &BindingParameters{}
}

func (d *dbOnlyManager) GetEmptyBindingParameters() service.BindingParameters {
	return &BindingParameters{}
}

func (a *allInOneManager) GetEmptyBindingDetails() service.BindingDetails {
	return &mysqlBindingDetails{}
}

func (v *vmOnlyManager) GetEmptyBindingDetails() service.BindingDetails {
	return &mysqlBindingDetails{}
}

func (d *dbOnlyManager) GetEmptyBindingDetails() service.BindingDetails {
	return &mysqlBindingDetails{}
}

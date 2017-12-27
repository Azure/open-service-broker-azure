package sqldb

import "github.com/Azure/open-service-broker-azure/pkg/service"

// ServerProvisioningParams encapsulates MSSQL-server specific provisioning
//options
type ServerProvisioningParams struct {
	FirewallIPStart string `json:"firewallStartIPAddress"`
	FirewallIPEnd   string `json:"firewallEndIPAddress"`
}

// DBProvisioningParams encapsulates MSSQL-specific provisioning options
type DBProvisioningParams struct {
}

type mssqlAllInOneInstanceDetails struct {
	ARMDeploymentName          string `json:"armDeployment"`
	ServerName                 string `json:"server"`
	AdministratorLogin         string `json:"administratorLogin"`
	AdministratorLoginPassword string `json:"administratorLoginPassword"`
	DatabaseName               string `json:"database"`
	FullyQualifiedDomainName   string `json:"fullyQualifiedDomainName"`
}

type mssqlVMOnlyInstanceDetails struct {
	ARMDeploymentName          string `json:"armDeployment"`
	ServerName                 string `json:"server"`
	AdministratorLogin         string `json:"administratorLogin"`
	AdministratorLoginPassword string `json:"administratorLoginPassword"`
	FullyQualifiedDomainName   string `json:"fullyQualifiedDomainName"`
}

type mssqlDBOnlyInstanceDetails struct {
	ARMDeploymentName string `json:"armDeployment"`
	ServerName        string `json:"server"`
	DatabaseName      string `json:"database"`
}

// UpdatingParameters encapsulates MSSQL-specific updating options
type UpdatingParameters struct {
}

// BindingParameters encapsulates MSSQL-specific binding options
type BindingParameters struct {
}

type mssqlBindingDetails struct {
	LoginName string `json:"loginName"`
	Password  string `json:"password"`
}

// Credentials encapsulates MSSQL-specific coonection details and credentials.
type Credentials struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Database string `json:"database"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// ServerConfig represents all configuration details needed for connecting to
// an Azure SQL Server.
type ServerConfig struct {
	ServerName                 string `json:"serverName"`
	ResourceGroupName          string `json:"resourceGroup"`
	Location                   string `json:"location"`
	AdministratorLogin         string `json:"administratorLogin"`
	AdministratorLoginPassword string `json:"administratorLoginPassword"`
}

// Config contains only a map of ServerConfig
type Config struct {
	Servers map[string]ServerConfig
}

func (
	a *allInOneManger,
) GetEmptyProvisioningParameters() service.ProvisioningParameters {
	return &ServerProvisioningParams{}
}

func (
	v *vmOnlyManager,
) GetEmptyProvisioningParameters() service.ProvisioningParameters {
	return &ServerProvisioningParams{}
}

func (
	d *dbOnlyManager,
) GetEmptyProvisioningParameters() service.ProvisioningParameters {
	return &DBProvisioningParams{}
}

func (
	a *allInOneManger,
) GetEmptyUpdatingParameters() service.UpdatingParameters {
	return &UpdatingParameters{}
}

func (
	v *vmOnlyManager,
) GetEmptyUpdatingParameters() service.UpdatingParameters {
	return &UpdatingParameters{}
}

func (
	d *dbOnlyManager,
) GetEmptyUpdatingParameters() service.UpdatingParameters {
	return &UpdatingParameters{}
}

func (
	a *allInOneManger,
) GetEmptyInstanceDetails() service.InstanceDetails {
	return &mssqlAllInOneInstanceDetails{}
}

func (
	v *vmOnlyManager,
) GetEmptyInstanceDetails() service.InstanceDetails {
	return &mssqlVMOnlyInstanceDetails{}
}

func (
	d *dbOnlyManager,
) GetEmptyInstanceDetails() service.InstanceDetails {
	return &mssqlDBOnlyInstanceDetails{}
}

func (
	a *allInOneManger,
) GetEmptyBindingParameters() service.BindingParameters {
	return &BindingParameters{}
}

func (
	v *vmOnlyManager,
) GetEmptyBindingParameters() service.BindingParameters {
	return &BindingParameters{}
}

func (
	d *dbOnlyManager,
) GetEmptyBindingParameters() service.BindingParameters {
	return &BindingParameters{}
}

func (
	a *allInOneManger,
) GetEmptyBindingDetails() service.BindingDetails {
	return &mssqlBindingDetails{}
}

func (
	v *vmOnlyManager,
) GetEmptyBindingDetails() service.BindingDetails {
	return &mssqlBindingDetails{}
}

func (
	d *dbOnlyManager,
) GetEmptyBindingDetails() service.BindingDetails {
	return &mssqlBindingDetails{}
}

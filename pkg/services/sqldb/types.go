package sqldb

import "github.com/Azure/open-service-broker-azure/pkg/service"

// ServerProvisioningParameters encapsulates MSSQL-specific provisioning options
// for provisioning involving server creation
type ServerProvisioningParameters struct {
	ServerName      string `json:"server"`
	FirewallIPStart string `json:"firewallStartIPAddress"`
	FirewallIPEnd   string `json:"firewallEndIPAddress"`
}

// DatabaseProvisioningParameters encapsulates MSSQL-specific provisioning
// options for provisioning involving db only provisioning
type DatabaseProvisioningParameters struct {
	ServerName string `json:"server"`
}

type mssqlServerOnlyInstanceDetails struct {
	ARMDeploymentName          string `json:"armDeployment"`
	ServerName                 string `json:"server"`
	AdministratorLogin         string `json:"administratorLogin"`
	AdministratorLoginPassword string `json:"administratorLoginPassword"`
	FullyQualifiedDomainName   string `json:"fullyQualifiedDomainName"`
}

type mssqlDBOnlyInstanceDetails struct {
	ARMDeploymentName string `json:"armDeployment"`
	DatabaseName      string `json:"database"`
}

type mssqlAllInOneInstanceDetails struct {
	ARMDeploymentName          string `json:"armDeployment"`
	ServerName                 string `json:"server"`
	AdministratorLogin         string `json:"administratorLogin"`
	AdministratorLoginPassword string `json:"administratorLoginPassword"`
	DatabaseName               string `json:"database"`
	FullyQualifiedDomainName   string `json:"fullyQualifiedDomainName"`
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

func (a *allServiceManager) GetEmptyProvisioningParameters() service.ProvisioningParameters {
	return &ServerProvisioningParameters{}
}

func (s *vmServiceManager) GetEmptyProvisioningParameters() service.ProvisioningParameters {
	return &ServerProvisioningParameters{}
}

func (d *dbServiceManager) GetEmptyProvisioningParameters() service.ProvisioningParameters {
	return &DatabaseProvisioningParameters{}
}

func (a *allServiceManager) GetEmptyUpdatingParameters() service.UpdatingParameters {
	return &UpdatingParameters{}
}

func (s *vmServiceManager) GetEmptyUpdatingParameters() service.UpdatingParameters {
	return &UpdatingParameters{}
}

func (d *dbServiceManager) GetEmptyUpdatingParameters() service.UpdatingParameters {
	return &UpdatingParameters{}
}

func (a *allServiceManager) GetEmptyInstanceDetails() service.InstanceDetails {
	return &mssqlAllInOneInstanceDetails{}
}

func (d *dbServiceManager) GetEmptyInstanceDetails() service.InstanceDetails {
	return &mssqlDBOnlyInstanceDetails{}
}

func (s *vmServiceManager) GetEmptyInstanceDetails() service.InstanceDetails {
	return &mssqlServerOnlyInstanceDetails{}
}

func (a *allServiceManager) GetEmptyBindingParameters() service.BindingParameters {
	return &BindingParameters{}
}

func (a *allServiceManager) GetEmptyBindingDetails() service.BindingDetails {
	return &mssqlBindingDetails{}
}

func (a *allServiceManager) GetEmptyCredentials() service.Credentials {
	return &Credentials{}
}

func (s *vmServiceManager) GetEmptyBindingParameters() service.BindingParameters {
	return &BindingParameters{}
}

func (s *vmServiceManager) GetEmptyBindingDetails() service.BindingDetails {
	return &mssqlBindingDetails{}
}

func (d *dbServiceManager) GetEmptyBindingParameters() service.BindingParameters {
	return &BindingParameters{}
}

func (d *dbServiceManager) GetEmptyBindingDetails() service.BindingDetails {
	return &mssqlBindingDetails{}
}

func (d *dbServiceManager) GetEmptyCredentials() service.Credentials {
	return &Credentials{}
}

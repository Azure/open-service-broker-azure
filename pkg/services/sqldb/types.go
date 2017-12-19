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

type mssqlServerOnlyProvisioningContext struct {
	ARMDeploymentName          string `json:"armDeployment"`
	ServerName                 string `json:"server"`
	AdministratorLogin         string `json:"administratorLogin"`
	AdministratorLoginPassword string `json:"administratorLoginPassword"`
	FullyQualifiedDomainName   string `json:"fullyQualifiedDomainName"`
}

type mssqlDBOnlyProvisioningContext struct {
	ARMDeploymentName string `json:"armDeployment"`
	DatabaseName      string `json:"database"`
}

type mssqlAllInOneProvisioningContext struct {
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

type mssqlBindingContext struct {
	LoginName string `json:"loginName"`
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

func (a *allServiceManager) GetEmptyProvisioningContext() service.ProvisioningContext {
	return &mssqlAllInOneProvisioningContext{}
}

func (d *dbServiceManager) GetEmptyProvisioningContext() service.ProvisioningContext {
	return &mssqlDBOnlyProvisioningContext{}
}

func (s *vmServiceManager) GetEmptyProvisioningContext() service.ProvisioningContext {
	return &mssqlServerOnlyProvisioningContext{}
}

func (a *allServiceManager) GetEmptyBindingParameters() service.BindingParameters {
	return &BindingParameters{}
}

func (a *allServiceManager) GetEmptyBindingContext() service.BindingContext {
	return &mssqlBindingContext{}
}

func (a *allServiceManager) GetEmptyCredentials() service.Credentials {
	return &Credentials{}
}

func (s *vmServiceManager) GetEmptyBindingParameters() service.BindingParameters {
	return &BindingParameters{}
}

func (s *vmServiceManager) GetEmptyBindingContext() service.BindingContext {
	return &mssqlBindingContext{}
}

func (s *vmServiceManager) GetEmptyCredentials() service.Credentials {
	return &Credentials{}
}

func (d *dbServiceManager) GetEmptyBindingParameters() service.BindingParameters {
	return &BindingParameters{}
}

func (d *dbServiceManager) GetEmptyBindingContext() service.BindingContext {
	return &mssqlBindingContext{}
}

func (d *dbSe) GetEmptyCredentials() service.Credentials {
	return &Credentials{}
}

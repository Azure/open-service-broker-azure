package sqldb

import "github.com/Azure/open-service-broker-azure/pkg/service"

// ProvisioningParameters encapsulates MSSQL-specific provisioning options
type ProvisioningParameters struct {
	ServerName      string `json:"server"`
	FirewallIPStart string `json:"firewallStartIPAddress"`
	FirewallIPEnd   string `json:"firewallEndIPAddress"`
}

type mssqlInstanceDetails struct {
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

func (
	s *serviceManager,
) GetEmptyProvisioningParameters() service.ProvisioningParameters {
	return &ProvisioningParameters{}
}

func (
	s *serviceManager,
) GetEmptyUpdatingParameters() service.UpdatingParameters {
	return &UpdatingParameters{}
}

func (
	s *serviceManager,
) GetEmptyInstanceDetails() service.InstanceDetails {
	return &mssqlInstanceDetails{}
}

func (s *serviceManager) GetEmptyBindingParameters() service.BindingParameters {
	return &BindingParameters{}
}

func (s *serviceManager) GetEmptyBindingDetails() service.BindingDetails {
	return &mssqlBindingDetails{}
}

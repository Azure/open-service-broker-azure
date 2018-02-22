package sqldb

// ServerProvisioningParams encapsulates non-sensitive MSSQL-server specific
// provisioning options
type ServerProvisioningParams struct {
	FirewallIPStart string `json:"firewallStartIPAddress"`
	FirewallIPEnd   string `json:"firewallEndIPAddress"`
}

// SecureServerProvisioningParams encapsulates sensitive MSSQL-server specific
// provisioning options
type SecureServerProvisioningParams struct{}

type serverInstanceDetails struct {
	FullyQualifiedDomainName string `json:"fullyQualifiedDomainName"`
	ServerName               string `json:"server"`
	AdministratorLogin       string `json:"administratorLogin"`
}

// UpdatingParameters encapsulates MSSQL-specific updating options
type UpdatingParameters struct {
}

// BindingParameters encapsulates non-sensitive MSSQL-specific binding options
type BindingParameters struct {
}

// SecureBindingParameters encapsulates sensitive MSSQL-specific binding options
type SecureBindingParameters struct {
}

type mssqlBindingDetails struct {
	LoginName string `json:"loginName"`
}

type mssqlSecureBindingDetails struct {
	Password string `json:"password"`
}

// Config contains only a map of ServerConfig
type Config struct {
	Servers map[string]ServerConfig
}

// Credentials encapsulates MSSQL-specific coonection details and credentials.
type Credentials struct {
	Host     string   `json:"host"`
	Port     int      `json:"port"`
	Database string   `json:"database"`
	Username string   `json:"username"`
	Password string   `json:"password"`
	URI      string   `json:"uri"`
	Tags     []string `json:"tags"`
	JDBC     string   `json:"jdbcUrl"`
	Encrypt  bool     `json:"encrypt"`
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

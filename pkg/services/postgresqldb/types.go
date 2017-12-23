package postgresqldb

import "github.com/Azure/open-service-broker-azure/pkg/service"

// ProvisioningParameters encapsulates PostgreSQL-specific provisioning options
type ProvisioningParameters struct {
	SSLEnforcement  string   `json:"sslEnforcement"`
	Extensions      []string `json:"extensions"`
	FirewallIPStart string   `json:"firewallStartIPAddress"`
	FirewallIPEnd   string   `json:"firewallEndIPAddress"`
}

type postgresqlInstanceDetails struct {
	ARMDeploymentName          string `json:"armDeployment"`
	ServerName                 string `json:"server"`
	AdministratorLoginPassword string `json:"administratorLoginPassword"`
	DatabaseName               string `json:"database"`
	FullyQualifiedDomainName   string `json:"fullyQualifiedDomainName"`
	EnforceSSL                 bool   `json:"enforceSSL"`
}

// UpdatingParameters encapsulates PostgreSQL-specific updating options
type UpdatingParameters struct {
}

// BindingParameters encapsulates PostgreSQL-specific binding options
type BindingParameters struct {
}

type postgresqlBindingDetails struct {
	LoginName string `json:"loginName"`
	Password  string `json:"password"`
}

// Credentials encapsulates PostgreSQL-specific coonection details and
// credentials.
type Credentials struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Database string `json:"database"`
	Username string `json:"username"`
	Password string `json:"password"`
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
	return &postgresqlInstanceDetails{}
}

func (s *serviceManager) GetEmptyBindingParameters() service.BindingParameters {
	return &BindingParameters{}
}

func (s *serviceManager) GetEmptyBindingDetails() service.BindingDetails {
	return &postgresqlBindingDetails{}
}

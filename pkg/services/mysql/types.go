package mysql

import "github.com/Azure/azure-service-broker/pkg/service"

// ProvisioningParameters encapsulates MySQL-specific provisioning options
type ProvisioningParameters struct {
	SSLEnforcement string `json:"sslEnforcement"`
}

type mysqlProvisioningContext struct {
	ARMDeploymentName          string `json:"armDeployment"`
	ServerName                 string `json:"server"`
	AdministratorLoginPassword string `json:"administratorLoginPassword"`
	DatabaseName               string `json:"database"`
	FullyQualifiedDomainName   string `json:"fullyQualifiedDomainName"`
	EnforceSSL                 bool   `json:"enforceSSL"`
}

// UpdatingParameters encapsulates MySQL-specific updating options
type UpdatingParameters struct {
}

// BindingParameters encapsulates MySQL-specific binding options
type BindingParameters struct {
}

type mysqlBindingContext struct {
	LoginName string `json:"loginName"`
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
	m *module,
) GetEmptyProvisioningParameters() service.ProvisioningParameters {
	return &ProvisioningParameters{}
}

func (
	m *module,
) GetEmptyUpdatingParameters() service.UpdatingParameters {
	return &UpdatingParameters{}
}

func (m *module) GetEmptyProvisioningContext() service.ProvisioningContext {
	return &mysqlProvisioningContext{}
}

func (m *module) GetEmptyBindingParameters() service.BindingParameters {
	return &BindingParameters{}
}

func (m *module) GetEmptyBindingContext() service.BindingContext {
	return &mysqlBindingContext{}
}

func (m *module) GetEmptyCredentials() service.Credentials {
	return &Credentials{}
}

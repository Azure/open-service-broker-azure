package postgresql

import "github.com/Azure/azure-service-broker/pkg/service"

// ProvisioningParameters encapsulates PostgreSQL-specific provisioning options
type ProvisioningParameters struct {
	Location      string            `json:"location"`
	ResourceGroup string            `json:"resourceGroup"`
	Tags          map[string]string `json:"tags"`
}

type postgresqlProvisioningContext struct {
	ResourceGroupName          string `json:"resourceGroup"`
	ARMDeploymentName          string `json:"armDeployment"`
	ServerName                 string `json:"server"`
	AdministratorLoginPassword string `json:"administratorLoginPassword"`
	DatabaseName               string `json:"database"`
	FullyQualifiedDomainName   string `json:"fullyQualifiedDomainName"`
}

// BindingParameters encapsulates PostgreSQL-specific binding options
type BindingParameters struct {
}

type postgresqlBindingContext struct {
	LoginName string `json:"loginName"`
}

type postgresqlCredentials struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Database string `json:"database"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func (p *postgresqlProvisioningContext) GetResourceGroupName() string {
	return p.ResourceGroupName
}

func (
	m *module,
) GetEmptyProvisioningParameters() service.ProvisioningParameters {
	return &ProvisioningParameters{}
}

func (m *module) GetEmptyProvisioningContext() service.ProvisioningContext {
	return &postgresqlProvisioningContext{}
}

func (m *module) GetEmptyBindingParameters() service.BindingParameters {
	return &BindingParameters{}
}

func (m *module) GetEmptyBindingContext() service.BindingContext {
	return &postgresqlBindingContext{}
}

func (m *module) GetEmptyCredentials() service.Credentials {
	return &postgresqlCredentials{}
}

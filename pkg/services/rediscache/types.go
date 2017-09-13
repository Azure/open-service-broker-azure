package rediscache

import "github.com/Azure/azure-service-broker/pkg/service"

// ProvisioningParameters encapsulates Redis-specific provisioning options
type ProvisioningParameters struct {
	Location      string            `json:"location"`
	ResourceGroup string            `json:"resourceGroup"`
	Tags          map[string]string `json:"tags"`
}

type redisProvisioningContext struct {
	ResourceGroupName        string `json:"resourceGroup"`
	ARMDeploymentName        string `json:"armDeployment"`
	ServerName               string `json:"server"`
	PrimaryKey               string `json:"primaryKey"`
	FullyQualifiedDomainName string `json:"fullyQualifiedDomainName"`
}

// BindingParameters encapsulates Redis-specific binding options
type BindingParameters struct {
}

type redisBindingContext struct {
}

type redisCredentials struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Password string `json:"password"`
}

func (r *redisProvisioningContext) GetResourceGroupName() string {
	return r.ResourceGroupName
}

func (
	m *module,
) GetEmptyProvisioningParameters() service.ProvisioningParameters {
	return &ProvisioningParameters{}
}

func (m *module) GetEmptyProvisioningContext() service.ProvisioningContext {
	return &redisProvisioningContext{}
}

func (m *module) GetEmptyBindingParameters() service.BindingParameters {
	return &BindingParameters{}
}

func (m *module) GetEmptyBindingContext() service.BindingContext {
	return &redisBindingContext{}
}

func (m *module) GetEmptyCredentials() service.Credentials {
	return &redisCredentials{}
}

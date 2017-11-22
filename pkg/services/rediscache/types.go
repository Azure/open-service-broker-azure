package rediscache

import "github.com/Azure/azure-service-broker/pkg/service"

// ProvisioningParameters encapsulates Redis-specific provisioning options
type ProvisioningParameters struct{}

type redisProvisioningContext struct {
	ARMDeploymentName        string `json:"armDeployment"`
	ServerName               string `json:"server"`
	PrimaryKey               string `json:"primaryKey"`
	FullyQualifiedDomainName string `json:"fullyQualifiedDomainName"`
}

// UpdatingParameters encapsulates Redis-specific updating options
type UpdatingParameters struct {
}

// BindingParameters encapsulates Redis-specific binding options
type BindingParameters struct {
}

type redisBindingContext struct {
}

// Credentials encapsulates Redis-specific coonection details and credentials.
type Credentials struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
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
	return &redisProvisioningContext{}
}

func (m *module) GetEmptyBindingParameters() service.BindingParameters {
	return &BindingParameters{}
}

func (m *module) GetEmptyBindingContext() service.BindingContext {
	return &redisBindingContext{}
}

func (m *module) GetEmptyCredentials() service.Credentials {
	return &Credentials{}
}

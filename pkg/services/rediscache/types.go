package rediscache

import "github.com/Azure/open-service-broker-azure/pkg/service"

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
) GetEmptyProvisioningContext() service.ProvisioningContext {
	return &redisProvisioningContext{}
}

func (s *serviceManager) GetEmptyBindingParameters() service.BindingParameters {
	return &BindingParameters{}
}

func (s *serviceManager) GetEmptyBindingContext() service.BindingContext {
	return &redisBindingContext{}
}

func (s *serviceManager) GetEmptyCredentials() service.Credentials {
	return &Credentials{}
}

package rediscache

import "github.com/Azure/open-service-broker-azure/pkg/service"

// ProvisioningParameters encapsulates non-sensitive Redis-specific provisioning
// options
type ProvisioningParameters struct{}

// SecureProvisioningParameters encapsulates sensitive Redis-specific
// provisioning options
type SecureProvisioningParameters struct{}

type redisInstanceDetails struct {
	ARMDeploymentName        string `json:"armDeployment"`
	ServerName               string `json:"server"`
	FullyQualifiedDomainName string `json:"fullyQualifiedDomainName"`
}

type redisSecureInstanceDetails struct {
	PrimaryKey string `json:"primaryKey"`
}

// UpdatingParameters encapsulates Redis-specific updating options
type UpdatingParameters struct {
}

// BindingParameters encapsulates Redis-specific binding options
type BindingParameters struct {
}

type redisBindingDetails struct {
}

type redisSecureBindingDetails struct {
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
) GetEmptySecureProvisioningParameters() service.SecureProvisioningParameters {
	return &SecureProvisioningParameters{}
}

func (
	s *serviceManager,
) GetEmptyUpdatingParameters() service.UpdatingParameters {
	return &UpdatingParameters{}
}

func (
	s *serviceManager,
) GetEmptyInstanceDetails() service.InstanceDetails {
	return &redisInstanceDetails{}
}

func (
	s *serviceManager,
) GetEmptySecureInstanceDetails() service.SecureInstanceDetails {
	return &redisSecureInstanceDetails{}
}

func (s *serviceManager) GetEmptyBindingParameters() service.BindingParameters {
	return &BindingParameters{}
}

func (s *serviceManager) GetEmptyBindingDetails() service.BindingDetails {
	return &redisBindingDetails{}
}

func (
	s *serviceManager,
) GetEmptySecureBindingDetails() service.SecureBindingDetails {
	return &redisSecureBindingDetails{}
}

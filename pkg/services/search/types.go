package search

import "github.com/Azure/open-service-broker-azure/pkg/service"

// ProvisioningParameters encapsulates non-sensitive Azure Search-specific
// provisioning options
type ProvisioningParameters struct{}

// SecureProvisioningParameters encapsulates sensitive Azure Search-specific
// provisioning options
type SecureProvisioningParameters struct{}

type searchInstanceDetails struct {
	ARMDeploymentName string `json:"armDeployment"`
	ServiceName       string `json:"serviceName"`
}

type searchSecureInstanceDetails struct {
	APIKey string `json:"apiKey"`
}

// UpdatingParameters encapsulates search-specific updating options
type UpdatingParameters struct {
}

// BindingParameters encapsulates non-sensitive Azure Search-specific binding
// options
type BindingParameters struct {
}

// SecureBindingParameters encapsulates sensitive Azure Search-specific binding
// options
type SecureBindingParameters struct {
}

type searchBindingDetails struct {
}

type searchSecureBindingDetails struct {
}

type searchCredentials struct {
	ServiceName string `json:"serviceName"`
	APIKey      string `json:"apiKey"`
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
	return &searchInstanceDetails{}
}

func (
	s *serviceManager,
) GetEmptySecureInstanceDetails() service.SecureInstanceDetails {
	return &searchSecureInstanceDetails{}
}

func (s *serviceManager) GetEmptyBindingParameters() service.BindingParameters {
	return &BindingParameters{}
}

func (
	s *serviceManager,
) GetEmptySecureBindingParameters() service.SecureBindingParameters {
	return &SecureBindingParameters{}
}

func (s *serviceManager) GetEmptyBindingDetails() service.BindingDetails {
	return &searchBindingDetails{}
}

func (
	s *serviceManager,
) GetEmptySecureBindingDetails() service.SecureBindingDetails {
	return &searchSecureBindingDetails{}
}

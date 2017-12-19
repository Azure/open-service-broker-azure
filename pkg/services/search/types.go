package search

import "github.com/Azure/open-service-broker-azure/pkg/service"

// ProvisioningParameters encapsulates
// Azure Search-specific provisioning options
type ProvisioningParameters struct{}

type searchInstanceDetails struct {
	ARMDeploymentName string `json:"armDeployment"`
	ServiceName       string `json:"serviceName"`
	APIKey            string `json:"apiKey"`
}

// UpdatingParameters encapsulates search-specific updating options
type UpdatingParameters struct {
}

// BindingParameters encapsulates Azure Search-specific binding options
type BindingParameters struct {
}

type searchBindingDetails struct {
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
) GetEmptyUpdatingParameters() service.UpdatingParameters {
	return &UpdatingParameters{}
}

func (
	s *serviceManager,
) GetEmptyInstanceDetails() service.InstanceDetails {
	return &searchInstanceDetails{}
}

func (s *serviceManager) GetEmptyBindingParameters() service.BindingParameters {
	return &BindingParameters{}
}

func (s *serviceManager) GetEmptyBindingDetails() service.BindingDetails {
	return &searchBindingDetails{}
}

func (s *serviceManager) GetEmptyCredentials() service.Credentials {
	return &searchCredentials{}
}

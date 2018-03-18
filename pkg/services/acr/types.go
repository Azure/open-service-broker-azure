package acr

import "github.com/Azure/open-service-broker-azure/pkg/service"

// ProvisioningParameters encapsulates
// Azure acr-specific provisioning options
type ProvisioningParameters struct{}

type acrProvisioningContext struct {
	ARMDeploymentName string `json:"armDeployment"`
	RegistryName      string `json:"registryName"`
	AdminUserEnabled  bool   `json:"adminUserEnabled"`
}

// UpdatingParameters encapsulates acr-specific updating options
type UpdatingParameters struct {
}

// BindingParameters encapsulates Azure acr-specific binding options
type BindingParameters struct {
}

type acrBindingContext struct {
}

type acrCredentials struct {
	RegistryName string `json:"registryName"`
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
	return &acrProvisioningContext{}
}

func (s *serviceManager) GetEmptyBindingParameters() service.BindingParameters {
	return &BindingParameters{}
}

func (s *serviceManager) GetEmptyBindingContext() service.BindingContext {
	return &acrBindingContext{}
}

func (s *serviceManager) GetEmptyCredentials() service.Credentials {
	return &acrCredentials{}
}

package search

import "github.com/Azure/azure-service-broker/pkg/service"

// ProvisioningParameters encapsulates
// Azure Search-specific provisioning options
type ProvisioningParameters struct{}

type searchProvisioningContext struct {
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

type searchBindingContext struct {
}

type searchCredentials struct {
	ServiceName string `json:"serviceName"`
	APIKey      string `json:"apiKey"`
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
	return &searchProvisioningContext{}
}

func (m *module) GetEmptyBindingParameters() service.BindingParameters {
	return &BindingParameters{}
}

func (m *module) GetEmptyBindingContext() service.BindingContext {
	return &searchBindingContext{}
}

func (m *module) GetEmptyCredentials() service.Credentials {
	return &searchCredentials{}
}

package aci

import "github.com/Azure/azure-service-broker/pkg/service"

// ProvisioningParameters encapsulates aci-specific provisioning options
type ProvisioningParameters struct {
	Location      string            `json:"location"`
	ResourceGroup string            `json:"resourceGroup"`
	Tags          map[string]string `json:"tags"`
	ImageName     string            `json:"image"`
	NumberCores   int               `json:"cpuCores"`
	Memory        float64           `json:"memoryInGb"`
	Port          int               `json:"port"`
}

type aciProvisioningContext struct {
	ResourceGroupName string `json:"resourceGroup"`
	ARMDeploymentName string `json:"armDeployment"`
	ContainerName     string `json:"name"`
	IPAddress         string `json:"containerIPv4Address"`
}

// UpdatingParameters encapsulates aci-specific updating options
type UpdatingParameters struct {
}

// BindingParameters encapsulates aci-specific binding options
type BindingParameters struct {
}

type aciBindingContext struct {
}

type aciCredentials struct {
	IPAddress string `json:"containerIPv4Address"`
}

func (
	m *module,
) GetEmptyProvisioningParameters() service.ProvisioningParameters {
	return &ProvisioningParameters{}
}

// SetResourceGroup sets the name of the resource group into which service
// instances will be deployed
func (p *ProvisioningParameters) SetResourceGroup(resourceGroup string) {
	p.ResourceGroup = resourceGroup
}

func (
	m *module,
) GetEmptyUpdatingParameters() service.UpdatingParameters {
	return &UpdatingParameters{}
}

func (m *module) GetEmptyProvisioningContext() service.ProvisioningContext {
	return &aciProvisioningContext{}
}

func (m *module) GetEmptyBindingParameters() service.BindingParameters {
	return &BindingParameters{}
}

func (m *module) GetEmptyBindingContext() service.BindingContext {
	return &aciBindingContext{}
}

func (m *module) GetEmptyCredentials() service.Credentials {
	return &aciCredentials{}
}

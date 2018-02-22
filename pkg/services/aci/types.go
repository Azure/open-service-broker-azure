package aci

import "github.com/Azure/open-service-broker-azure/pkg/service"

// ProvisioningParameters encapsulates non-sensitive aci-specific provisioning
// options
type ProvisioningParameters struct {
	ImageName   string  `json:"image"`
	NumberCores int     `json:"cpuCores"`
	Memory      float64 `json:"memoryInGb"`
	Ports       []int   `json:"ports"`
}

// SecureProvisioningParameters encapsulates sensitive aci-specific provisioning
// options
type SecureProvisioningParameters struct{}

type aciInstanceDetails struct {
	ARMDeploymentName string `json:"armDeployment"`
	ContainerName     string `json:"name"`
	PublicIPv4Address string `json:"publicIPv4Address"`
}

type aciSecureInstanceDetails struct{}

// BindingParameters encapsulates non-sensitive aci-specific binding options
type BindingParameters struct {
}

// SecureBindingParameters encapsulates sensitive aci-specific binding options
type SecureBindingParameters struct {
}

type aciBindingDetails struct {
}

type aciSecureBindingDetails struct {
}

type aciCredentials struct {
	PublicIPv4Address string `json:"publicIPv4Address"`
}

func (
	s *serviceManager,
) GetEmptyProvisioningParameters() service.ProvisioningParameters {
	return &ProvisioningParameters{
		NumberCores: 1,
		Memory:      1.5,
		Ports:       make([]int, 0),
	}
}

func (
	s *serviceManager,
) GetEmptySecureProvisioningParameters() service.SecureProvisioningParameters {
	return &SecureProvisioningParameters{}
}

func (s *serviceManager) GetEmptyInstanceDetails() service.InstanceDetails {
	return &aciInstanceDetails{}
}

func (
	s *serviceManager,
) GetEmptySecureInstanceDetails() service.SecureInstanceDetails {
	return &aciSecureInstanceDetails{}
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
	return &aciBindingDetails{}
}

func (
	s *serviceManager,
) GetEmptySecureBindingDetails() service.SecureBindingDetails {
	return &aciSecureBindingDetails{}
}

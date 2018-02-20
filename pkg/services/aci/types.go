package aci

import "github.com/Azure/open-service-broker-azure/pkg/service"

// ProvisioningParameters encapsulates aci-specific provisioning options
type ProvisioningParameters struct {
	ImageName   string  `json:"image"`
	NumberCores int     `json:"cpuCores"`
	Memory      float64 `json:"memoryInGb"`
	Ports       []int   `json:"ports"`
}

type aciInstanceDetails struct {
	ARMDeploymentName string `json:"armDeployment"`
	ContainerName     string `json:"name"`
	PublicIPv4Address string `json:"publicIPv4Address"`
}

type aciSecureInstanceDetails struct{}

// UpdatingParameters encapsulates aci-specific updating options
type UpdatingParameters struct {
}

// BindingParameters encapsulates aci-specific binding options
type BindingParameters struct {
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
) GetEmptyUpdatingParameters() service.UpdatingParameters {
	return &UpdatingParameters{}
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

func (s *serviceManager) GetEmptyBindingDetails() service.BindingDetails {
	return &aciBindingDetails{}
}

func (
	s *serviceManager,
) GetEmptySecureBindingDetails() service.SecureBindingDetails {
	return &aciSecureBindingDetails{}
}

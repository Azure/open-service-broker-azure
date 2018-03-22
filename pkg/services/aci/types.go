package aci

import "github.com/Azure/open-service-broker-azure/pkg/service"

type provisioningParameters struct {
	ImageName   string  `json:"image"`
	NumberCores int     `json:"cpuCores"`
	Memory      float64 `json:"memoryInGb"`
	Ports       []int   `json:"ports"`
}

type instanceDetails struct {
	ARMDeploymentName string `json:"armDeployment"`
	ContainerName     string `json:"name"`
	PublicIPv4Address string `json:"publicIPv4Address"`
}

type credentials struct {
	PublicIPv4Address string `json:"publicIPv4Address"`
}

func (s *serviceManager) SplitProvisioningParameters(
	cpp service.CombinedProvisioningParameters,
) (
	service.ProvisioningParameters,
	service.SecureProvisioningParameters,
	error,
) {
	pp := provisioningParameters{
		NumberCores: 1,
		Memory:      1.5,
		Ports:       make([]int, 0),
	}
	if err := service.GetStructFromMap(cpp, &pp); err != nil {
		return nil, nil, err
	}
	ppMap, err := service.GetMapFromStruct(pp)
	if err != nil {
		return nil, nil, err
	}
	return ppMap, nil, nil
}

func (s *serviceManager) SplitBindingParameters(
	params service.CombinedBindingParameters,
) (
	service.BindingParameters,
	service.SecureBindingParameters,
	error,
) {
	return nil, nil, nil
}

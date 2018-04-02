package aci

import (
	"github.com/Azure/open-service-broker-azure/pkg/service"
)

type provisioningParameters struct {
	ImageName   string  `json:"image"`
	NumberCores int     `json:"cpuCores"`
	Memory      float64 `json:"memoryInGb"`
	Ports       []int   `json:"ports"`
}

func (
	s *serviceManager,
) getProvisionParametersSchema() map[string]service.ParameterSchema {

	p := map[string]service.ParameterSchema{}

	imageSchema := service.NewSimpleParameterSchema(
		"string",
		"The Docker image on which to base the container.",
	)
	imageSchema.SetRequired(true)
	p["image"] = imageSchema

	cpuCoreSchema := service.NewSimpleParameterSchema(
		"integer",
		"The number of virtual CPU cores requested "+
			"for the container.",
	)
	cpuCoreSchema.SetDefault(1)
	p["cpuCores"] = cpuCoreSchema

	memorySchema := service.NewSimpleParameterSchema(
		"integer",
		"Gigabytes of memory requested for the container. "+
			"Must be specified in increments of 0.10 GB.",
	)
	memorySchema.SetDefault(1.5)
	p["memoryInGb"] = memorySchema

	portsSchema := service.NewArrayParameterSchema(
		"The port(s) to open on the container. The container "+
			"will be assigned a public IP (v4) address if and only if one or "+
			"more ports are opened.",
		service.NewSimpleParameterSchema(
			"integer",
			"Port to open on container",
		),
	)
	portsSchema.SetRequired(true)
	p["ports"] = portsSchema
	return p
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

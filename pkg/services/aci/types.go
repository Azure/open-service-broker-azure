package aci

import "github.com/Azure/open-service-broker-azure/pkg/service"

type provisioningParameters struct {
	ImageName   string  `json:"image"`
	NumberCores int     `json:"cpuCores"`
	Memory      float64 `json:"memoryInGb"`
	Ports       []int   `json:"ports"`
}

// GetSchema generates the schema for instance provisioning parameters
func GetSchema() *service.ParametersSchema {

	p := service.GetCommonSchema()

	p.Properties["image"] = service.Parameter{
		Type:        "string",
		Description: "The Docker image on which to base the container.",
	}

	p.Properties["cpuCores"] = service.Parameter{
		Type: "integer",
		Description: "The number of virtual CPU cores requested " +
			"for the container.",
		Default: 1,
	}

	p.Properties["memoryInGb"] = service.Parameter{
		Type: "integer",
		Description: "Gigabytes of memory requested for the container. " +
			"Must be specified in increments of 0.10 GB.",
		Default: 1.5,
	}

	p.Properties["ports"] = service.Parameter{
		Type: "array",
		Description: "The port(s) to open on the container. The container " +
			"will be assigned a public IP (v4) address if and only if one or " +
			"more ports are opened.",
		Items: service.Parameter{
			Type: "integer",
		},
	}

	p.Required = []string{"image", "ports"}
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

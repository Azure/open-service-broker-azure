package aci

import (
	"github.com/Azure/open-service-broker-azure/pkg/ptr"
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
) getProvisionParametersSchema() service.InputParametersSchema {
	return service.InputParametersSchema{
		RequiredProperties: []string{"image"},
		PropertySchemas: map[string]service.PropertySchema{
			"image": &service.StringPropertySchema{
				Description: "The Docker image on which to base the container.",
			},
			"cpuCores": &service.IntPropertySchema{
				Description: "The number of virtual CPU cores requested " +
					"for the container.",
				DefaultValue: ptr.ToInt64(1),
			},
			"memoryInGb": &service.FloatPropertySchema{
				Description: "Gigabytes of memory requested for the container. " +
					"Must be specified in increments of 0.10 GB.",
				DefaultValue: ptr.ToFloat64(1.5),
				// krancour: Currently not supported because of floating point division
				// errors.
				// AllowedIncrement: ptr.ToFloat64(0.10),
				CustomPropertyValidator: memoryValidator,
			},
			"ports": &service.ArrayPropertySchema{
				Description: "The port(s) to open on the container." +
					"The container will be assigned a public IP (v4) address if" +
					" and only if one or more ports are opened.",
				ItemsSchema: &service.IntPropertySchema{
					Description: "Port to open on container",
				},
			},
		},
	}
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
	cpp map[string]interface{},
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

func memoryValidator(context string, value float64) error {
	value *= 10
	if float64(int64(value)) != value {
		return service.NewValidationError(
			context,
			"memory must be specified in increments of 0.10 GB",
		)
	}
	return nil
}

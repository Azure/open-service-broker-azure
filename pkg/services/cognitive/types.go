package cognitive

import (
	"github.com/Azure/open-service-broker-azure/pkg/service"
)

type instanceDetails struct {
	ARMDeploymentName string `json:"armDeployment"`
	CognitiveName     string `json:"name"`
	CognitiveKey      string `json:"cognitivekey"`
	Endpoint          string `json:"endpoint"`
}

type credentials struct {
	CognitiveKey string `json:"cognitivekey"`
	Endpoint     string `json:"endpoint"`
}

func (s *serviceManager) SplitProvisioningParameters(
	service.CombinedProvisioningParameters,
) (
	service.ProvisioningParameters,
	service.SecureProvisioningParameters,
	error,
) {
	return nil, nil, nil
}

func (s *serviceManager) SplitBindingParameters(
	service.CombinedBindingParameters,
) (service.BindingParameters, service.SecureBindingParameters, error) {
	return nil, nil, nil
}

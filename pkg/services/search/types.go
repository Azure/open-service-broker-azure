package search

import "github.com/Azure/open-service-broker-azure/pkg/service"

type instanceDetails struct {
	ARMDeploymentName string `json:"armDeployment"`
	ServiceName       string `json:"serviceName"`
}

type secureInstanceDetails struct {
	APIKey string `json:"apiKey"`
}

type credentials struct {
	ServiceName string `json:"serviceName"`
	APIKey      string `json:"apiKey"`
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

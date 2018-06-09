// +build experimental

package rediscache

import "github.com/Azure/open-service-broker-azure/pkg/service"

type instanceDetails struct {
	ARMDeploymentName        string `json:"armDeployment"`
	ServerName               string `json:"server"`
	FullyQualifiedDomainName string `json:"fullyQualifiedDomainName"`
}

type secureInstanceDetails struct {
	PrimaryKey string `json:"primaryKey"`
}

type credentials struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Password string `json:"password"`
	URI      string `json:"uri"`
}

func (s *serviceManager) SplitProvisioningParameters(
	map[string]interface{},
) (
	service.ProvisioningParameters,
	service.SecureProvisioningParameters,
	error,
) {
	return nil, nil, nil
}

func (s *serviceManager) SplitBindingParameters(
	map[string]interface{},
) (service.BindingParameters, service.SecureBindingParameters, error) {
	return nil, nil, nil
}

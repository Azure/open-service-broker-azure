package appinsights

import "github.com/Azure/open-service-broker-azure/pkg/service"

type instanceDetails struct {
	ARMDeploymentName  string               `json:"armDeployment"`
	InstrumentationKey service.SecureString `json:"instrumentationKey"`
}

type bindingDetails struct {
	AppID    string               `json:"appID"`
	APIKeyID string               `json:"APIKeyID"`
	APIKey   service.SecureString `json:"APIKey"`
}

func (s *serviceManager) GetEmptyInstanceDetails() service.InstanceDetails {
	return &instanceDetails{}
}

func (s *serviceManager) GetEmptyBindingDetails() service.BindingDetails {
	return nil
}

type credentials struct {
	InstrumentationKey string `json:"instrumentationKey"`
	AppID              string `json:"appID"`
	APIKey             string `json:"APIKey"`
}

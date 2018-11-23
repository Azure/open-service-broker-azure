package appinsights

import "github.com/Azure/open-service-broker-azure/pkg/service"

type instanceDetails struct {
	ARMDeploymentName  string               `json:"armDeployment"`
	AppInsightsName    string               `json:"appInsightsName"`
	InstrumentationKey service.SecureString `json:"instrumentationKey"`
}

func (s *serviceManager) GetEmptyInstanceDetails() service.InstanceDetails {
	return &instanceDetails{}
}

func (s *serviceManager) GetEmptyBindingDetails() service.BindingDetails {
	return nil
}

type credentials struct {
	InstrumentationKey string `json:"instrumentationKey"`
}

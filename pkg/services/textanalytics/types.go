package textanalytics

import "github.com/Azure/open-service-broker-azure/pkg/service"

type instanceDetails struct {
	ARMDeploymentName string `json:"armDeployment"`
	TextAnalyticsName string `json:"textAnalyticsName"`
	TextAnalyticsKey  string `json:"textAnalyticsKey"`
	Endpoint          string `json:"textAnalyticsEndpoint"`
}

type credentials struct {
	TextAnalyticsKey  string `json:"textAnalyticsKey"`
	Endpoint          string `json:"textAnalyticsEndpoint"`
	TextAnalyticsName string `json:"textAnalyticsName"`
}

func (s *serviceManager) GetEmptyInstanceDetails() service.InstanceDetails {
	return &instanceDetails{}
}

func (s *serviceManager) GetEmptyBindingDetails() service.BindingDetails {
	return nil
}

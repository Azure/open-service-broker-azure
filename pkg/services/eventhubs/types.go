package eventhubs

import "github.com/Azure/open-service-broker-azure/pkg/service"

type instanceDetails struct {
	ARMDeploymentName string               `json:"armDeployment"`
	EventHubName      string               `json:"eventHubName"`
	EventHubNamespace string               `json:"eventHubNamespace"`
	ConnectionString  service.SecureString `json:"connectionString"`
	PrimaryKey        service.SecureString `json:"primaryKey"`
}

func (s *serviceManager) GetEmptyInstanceDetails() service.InstanceDetails {
	return &instanceDetails{}
}

func (s *serviceManager) GetEmptyBindingDetails() service.BindingDetails {
	return nil
}

type credentials struct {
	ConnectionString  string               `json:"connectionString"`
	PrimaryKey        string               `json:"primaryKey"`
	EventHubNamespace string               `json:"eventHubNamespace"`
	EventHubName      string               `json:"eventHubName"`
}

package servicebus

import "github.com/Azure/open-service-broker-azure/pkg/service"

type instanceDetails struct {
	ARMDeploymentName       string               `json:"armDeployment"`
	ServiceBusNamespaceName string               `json:"serviceBusNamespaceName"`
	ConnectionString        service.SecureString `json:"connectionString"`
	PrimaryKey              service.SecureString `json:"primaryKey"`
}

func (s *serviceManager) GetEmptyInstanceDetails() service.InstanceDetails {
	return &instanceDetails{}
}

func (s *serviceManager) GetEmptyBindingDetails() service.BindingDetails {
	return nil
}

type credentials struct {
	ConnectionString string `json:"connectionString"`
	PrimaryKey       string `json:"primaryKey"`
}

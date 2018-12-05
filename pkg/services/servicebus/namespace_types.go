package servicebus

import "github.com/Azure/open-service-broker-azure/pkg/service"

type namespaceInstanceDetails struct {
	ARMDeploymentName       string               `json:"armDeployment"`
	ServiceBusNamespaceName string               `json:"serviceBusNamespaceName"`
	ConnectionString        service.SecureString `json:"connectionString"`
	PrimaryKey              service.SecureString `json:"primaryKey"`
}

func (nm *namespaceManager) GetEmptyInstanceDetails() service.InstanceDetails {
	return &namespaceInstanceDetails{}
}

func (nm *namespaceManager) GetEmptyBindingDetails() service.BindingDetails {
	return nil
}

type namespaceCredentials struct {
	ConnectionString string `json:"connectionString"`
	PrimaryKey       string `json:"primaryKey"`
}

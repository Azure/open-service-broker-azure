package servicebus

import "github.com/Azure/open-service-broker-azure/pkg/service"

type queueInstanceDetails struct {
	QueueName string `json:"queueName"`
}

func (qm *queueManager) GetEmptyInstanceDetails() service.InstanceDetails {
	return &queueInstanceDetails{}
}

func (qm *queueManager) GetEmptyBindingDetails() service.BindingDetails {
	return nil
}

type queueCredentials struct {
	ConnectionString string `json:"connectionString"`
	PrimaryKey       string `json:"primaryKey"`
	NamespaceName    string `json:"namespaceName"`
	QueueName        string `json:"queueName"`
}

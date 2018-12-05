package servicebus

import "github.com/Azure/open-service-broker-azure/pkg/service"

type queueInstanceDetails struct {
	ServiceBusQueueName string `json:"serviceBusQueueName"`
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
	QueueURL         string `json:"queueURL"`
}

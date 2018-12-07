package servicebus

import "github.com/Azure/open-service-broker-azure/pkg/service"

type topicInstanceDetails struct {
	ServiceBusTopicName string `json:"serviceBusTopicName"`
}

func (tm *topicManager) GetEmptyInstanceDetails() service.InstanceDetails {
	return &topicInstanceDetails{}
}

func (tm *topicManager) GetEmptyBindingDetails() service.BindingDetails {
	return nil
}

type topicCredentials struct {
	ConnectionString string `json:"connectionString"`
	PrimaryKey       string `json:"primaryKey"`
	TopicName        string `json:"topicName"`
	TopicURL         string `json:"topicURL"`
}

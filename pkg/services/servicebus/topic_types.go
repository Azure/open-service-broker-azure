package servicebus

import "github.com/Azure/open-service-broker-azure/pkg/service"

type topicInstanceDetails struct {
	TopicName string `json:"topicName"`
}

type topicBindingDetails struct {
	SubscriptionName string `json:"subscriptionName,omitempty"`
}

func (tm *topicManager) GetEmptyInstanceDetails() service.InstanceDetails {
	return &topicInstanceDetails{}
}

func (tm *topicManager) GetEmptyBindingDetails() service.BindingDetails {
	return &topicBindingDetails{}
}

type topicCredentials struct {
	ConnectionString string `json:"connectionString"`
	PrimaryKey       string `json:"primaryKey"`
	NamespaceName    string `json:"namespaceName"`
	TopicName        string `json:"topicName"`
	SubscriptionName string `json:"subscriptionName,omitempty"`
}

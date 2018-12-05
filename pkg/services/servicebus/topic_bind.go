package servicebus

import (
	"fmt"

	"github.com/Azure/open-service-broker-azure/pkg/service"
)

func (tm *topicManager) Bind(
	service.Instance,
	service.BindingParameters,
) (service.BindingDetails, error) {
	return nil, nil
}

func (tm *topicManager) GetCredentials(
	instance service.Instance,
	_ service.Binding,
) (service.Credentials, error) {
	pdt := instance.Parent.Details.(*namespaceInstanceDetails)
	dt := instance.Details.(*topicInstanceDetails)
	return topicCredentials{
		ConnectionString: string(pdt.ConnectionString),
		PrimaryKey:       string(pdt.PrimaryKey),
		TopicURL: fmt.Sprintf(
			"https://%s.servicebus.windows.net/%s",
			pdt.ServiceBusNamespaceName,
			dt.ServiceBusTopicName,
		),
	}, nil
}

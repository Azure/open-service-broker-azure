package servicebus

import (
	"fmt"

	"github.com/Azure/open-service-broker-azure/pkg/service"
)

func (qm *queueManager) Bind(
	service.Instance,
	service.BindingParameters,
) (service.BindingDetails, error) {
	return nil, nil
}

func (qm *queueManager) GetCredentials(
	instance service.Instance,
	_ service.Binding,
) (service.Credentials, error) {
	pdt := instance.Parent.Details.(*namespaceInstanceDetails)
	dt := instance.Details.(*queueInstanceDetails)
	return queueCredentials{
		ConnectionString: string(pdt.ConnectionString),
		PrimaryKey:       string(pdt.PrimaryKey),
		QueueName:        string(dt.ServiceBusQueueName),
		QueueURL: fmt.Sprintf(
			"https://%s.servicebus.windows.net/%s",
			pdt.ServiceBusNamespaceName,
			dt.ServiceBusQueueName,
		),
	}, nil
}

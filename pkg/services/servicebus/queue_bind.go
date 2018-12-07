package servicebus

import (
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
		NamespaceName:    string(pdt.NamespaceName),
		QueueName:        string(dt.QueueName),
	}, nil
}

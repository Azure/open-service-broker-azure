package servicebus

import (
	"context"

	"github.com/Azure/open-service-broker-azure/pkg/service"
)

func (tm *topicManager) Unbind(
	instance service.Instance,
	binding service.Binding,
) error {
	bt := binding.Details.(*topicBindingDetails)
	if bt.SubscriptionName == "" {
		return nil
	}

	ppp := instance.Parent.ProvisioningParameters
	pdt := instance.Parent.Details.(*namespaceInstanceDetails)
	dt := instance.Details.(*topicInstanceDetails)
	_, err := tm.subscriptionsClient.Delete(
		context.TODO(),
		ppp.GetString("resourceGroup"),
		pdt.ServiceBusNamespaceName,
		dt.ServiceBusTopicName,
		bt.SubscriptionName,
	)
	return err
}

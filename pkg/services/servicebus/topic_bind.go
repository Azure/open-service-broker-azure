package servicebus

import (
	"context"
	"fmt"

	servicebusSDK "github.com/Azure/azure-sdk-for-go/services/servicebus/mgmt/2017-04-01/servicebus" // nolint: lll
	"github.com/Azure/open-service-broker-azure/pkg/service"
	uuid "github.com/satori/go.uuid"
)

func (tm *topicManager) Bind(
	instance service.Instance,
	bindingParameters service.BindingParameters,
) (service.BindingDetails, error) {
	if bindingParameters.GetString("subscriptionNeeded") != "yes" {
		return &topicBindingDetails{}, nil
	}

	bt := &topicBindingDetails{
		SubscriptionName: uuid.NewV4().String(),
	}
	pdt := instance.Parent.Details.(*namespaceInstanceDetails)
	ppp := instance.Parent.ProvisioningParameters
	dt := instance.Details.(*topicInstanceDetails)

	// Using context.TODO() here because no parent context is passed in.
	if _, err := tm.subscriptionsClient.CreateOrUpdate(
		context.TODO(),
		ppp.GetString("resourceGroup"),
		pdt.NamespaceName,
		dt.TopicName,
		bt.SubscriptionName,
		servicebusSDK.SBSubscription{},
	); err != nil {
		return nil, fmt.Errorf("error creating subscription in binding: %s", err)
	}
	return bt, nil
}

func (tm *topicManager) GetCredentials(
	instance service.Instance,
	binding service.Binding,
) (service.Credentials, error) {
	pdt := instance.Parent.Details.(*namespaceInstanceDetails)
	dt := instance.Details.(*topicInstanceDetails)
	bt := binding.Details.(*topicBindingDetails)
	return topicCredentials{
		ConnectionString: string(pdt.ConnectionString),
		PrimaryKey:       string(pdt.PrimaryKey),
		NamespaceName:    pdt.NamespaceName,
		TopicName:        dt.TopicName,
		SubscriptionName: bt.SubscriptionName,
	}, nil
}

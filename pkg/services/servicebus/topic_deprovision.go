package servicebus

import (
	"context"
	"fmt"

	"github.com/Azure/open-service-broker-azure/pkg/service"
)

func (tm *topicManager) GetDeprovisioner(
	service.Plan,
) (service.Deprovisioner, error) {
	return service.NewDeprovisioner(
		service.NewDeprovisioningStep("deleteTopic", tm.deleteTopic),
	)
}

func (tm *topicManager) deleteTopic(
	ctx context.Context,
	instance service.Instance,
) (service.InstanceDetails, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	dt := instance.Details.(*topicInstanceDetails)
	pdt := instance.Parent.Details.(*namespaceInstanceDetails)
	ppp := instance.Parent.ProvisioningParameters
	_, err := tm.topicsClient.Delete(
		ctx,
		ppp.GetString("resourceGroup"),
		pdt.ServiceBusNamespaceName,
		dt.ServiceBusTopicName,
	)
	if err != nil {
		return nil, fmt.Errorf("error deleting service bus topic: %s", err)
	}

	return instance.Details, nil
}

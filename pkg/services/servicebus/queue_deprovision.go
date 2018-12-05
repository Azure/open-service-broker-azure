package servicebus

import (
	"context"
	"fmt"

	"github.com/Azure/open-service-broker-azure/pkg/service"
)

func (qm *queueManager) GetDeprovisioner(
	service.Plan,
) (service.Deprovisioner, error) {
	return service.NewDeprovisioner(
		service.NewDeprovisioningStep("deleteQueue", qm.deleteQueue),
	)
}

func (qm *queueManager) deleteQueue(
	ctx context.Context,
	instance service.Instance,
) (service.InstanceDetails, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	dt := instance.Details.(*queueInstanceDetails)
	pdt := instance.Parent.Details.(*namespaceInstanceDetails)
	ppp := instance.Parent.ProvisioningParameters
	_, err := qm.queuesClient.Delete(
		ctx,
		ppp.GetString("resourceGroup"),
		pdt.ServiceBusNamespaceName,
		dt.ServiceBusQueueName,
	)
	if err != nil {
		return nil, fmt.Errorf("error deleting service bus queue: %s", err)
	}

	return instance.Details, nil
}

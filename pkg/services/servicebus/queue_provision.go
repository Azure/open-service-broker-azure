package servicebus

import (
	"context"
	"fmt"
	"strings"

	"github.com/Azure/open-service-broker-azure/pkg/ptr"

	servicebusSDK "github.com/Azure/azure-sdk-for-go/services/servicebus/mgmt/2017-04-01/servicebus" // nolint: lll
	"github.com/Azure/open-service-broker-azure/pkg/service"
	uuid "github.com/satori/go.uuid"
)

func (qm *queueManager) GetProvisioner(
	service.Plan,
) (service.Provisioner, error) {
	return service.NewProvisioner(
		service.NewProvisioningStep("preProvision", qm.preProvision),
		service.NewProvisioningStep("createQueue", qm.createQueue),
	)
}

// nolint: lll
func (qm *queueManager) preProvision(
	ctx context.Context,
	instance service.Instance,
) (service.InstanceDetails, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	ppp := instance.Parent.ProvisioningParameters
	pdt := instance.Parent.Details.(*namespaceInstanceDetails)
	if queueName := instance.ProvisioningParameters.GetString("queueName"); queueName != "" {
		getResult, err := qm.queuesClient.Get(
			ctx,
			ppp.GetString("resourceGroup"),
			pdt.NamespaceName,
			queueName,
		)
		if getResult.Name != nil && err == nil {
			return nil, fmt.Errorf("queue with name '%s' already exists in the namespace", queueName)
		} else if !strings.Contains(err.Error(), "StatusCode=404") {
			return nil, fmt.Errorf("unexpected error creating queue: %s", err)
		} else {
			return &queueInstanceDetails{
				QueueName: queueName,
			}, nil
		}
	}
	return &queueInstanceDetails{
		QueueName: uuid.NewV4().String(),
	}, nil
}

func (qm *queueManager) createQueue(
	ctx context.Context,
	instance service.Instance,
) (service.InstanceDetails, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	ppp := instance.Parent.ProvisioningParameters
	pdt := instance.Parent.Details.(*namespaceInstanceDetails)
	dt := instance.Details.(*queueInstanceDetails)
	if _, err := qm.queuesClient.CreateOrUpdate(
		ctx,
		ppp.GetString("resourceGroup"),
		pdt.NamespaceName,
		dt.QueueName,
		qm.buildQueueInformation(instance),
	); err != nil {
		return nil, fmt.Errorf("error creating queue: %s", err)
	}
	return dt, nil
}

func (qm *queueManager) buildQueueInformation(
	instance service.Instance,
) servicebusSDK.SBQueue {
	pp := instance.ProvisioningParameters
	return servicebusSDK.SBQueue{
		SBQueueProperties: &servicebusSDK.SBQueueProperties{
			MaxSizeInMegabytes:       ptr.ToInt32(int32(pp.GetInt64("maxQueueSize"))),
			DefaultMessageTimeToLive: ptr.ToString(pp.GetString("messageTimeToLive")),
			LockDuration:             ptr.ToString(pp.GetString("lockDuration")),
		},
	}
}

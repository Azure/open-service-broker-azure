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

func (tm *topicManager) GetProvisioner(
	service.Plan,
) (service.Provisioner, error) {
	return service.NewProvisioner(
		service.NewProvisioningStep("preProvision", tm.preProvision),
		service.NewProvisioningStep("createTopic", tm.createTopic),
	)
}

// nolint: lll
func (tm *topicManager) preProvision(
	ctx context.Context,
	instance service.Instance,
) (service.InstanceDetails, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	ppp := instance.Parent.ProvisioningParameters
	pdt := instance.Parent.Details.(*namespaceInstanceDetails)
	if topicName := instance.ProvisioningParameters.GetString("topicName"); topicName != "" {
		getResult, err := tm.topicsClient.Get(
			ctx,
			ppp.GetString("resourceGroup"),
			pdt.NamespaceName,
			topicName,
		)
		if getResult.Name != nil && err == nil {
			return nil, fmt.Errorf("topic with name '%s' already exists in the namespace", topicName)
		} else if !strings.Contains(err.Error(), "StatusCode=404") {
			return nil, fmt.Errorf("unexpected error creating topic: %s", err)
		} else {
			return &topicInstanceDetails{
				TopicName: topicName,
			}, nil
		}
	}
	return &topicInstanceDetails{
		TopicName: uuid.NewV4().String(),
	}, nil
}

func (tm *topicManager) createTopic(
	ctx context.Context,
	instance service.Instance,
) (service.InstanceDetails, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	ppp := instance.Parent.ProvisioningParameters
	pdt := instance.Parent.Details.(*namespaceInstanceDetails)
	dt := instance.Details.(*topicInstanceDetails)
	if _, err := tm.topicsClient.CreateOrUpdate(
		ctx,
		ppp.GetString("resourceGroup"),
		pdt.NamespaceName,
		dt.TopicName,
		tm.buildTopicInformation(instance),
	); err != nil {
		return nil, fmt.Errorf("error creating topic: %s", err)
	}
	return dt, nil
}

func (tm *topicManager) buildTopicInformation(
	instance service.Instance,
) servicebusSDK.SBTopic {
	pp := instance.ProvisioningParameters
	return servicebusSDK.SBTopic{
		SBTopicProperties: &servicebusSDK.SBTopicProperties{
			MaxSizeInMegabytes:       ptr.ToInt32(int32(pp.GetInt64("maxTopicSize"))),
			DefaultMessageTimeToLive: ptr.ToString(pp.GetString("messageTimeToLive")),
		},
	}
}

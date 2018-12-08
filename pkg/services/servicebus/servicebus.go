package servicebus

import (
	servicebusSDK "github.com/Azure/azure-sdk-for-go/services/servicebus/mgmt/2017-04-01/servicebus" // nolint: lll
	"github.com/Azure/open-service-broker-azure/pkg/azure/arm"
	"github.com/Azure/open-service-broker-azure/pkg/service"
)

type module struct {
	namespaceManager *namespaceManager
	queueManager     *queueManager
	topicManager     *topicManager
}

type namespaceManager struct {
	armDeployer      arm.Deployer
	namespacesClient servicebusSDK.NamespacesClient
}

type queueManager struct {
	armDeployer  arm.Deployer
	queuesClient servicebusSDK.QueuesClient
}

type topicManager struct {
	armDeployer         arm.Deployer
	topicsClient        servicebusSDK.TopicsClient
	subscriptionsClient servicebusSDK.SubscriptionsClient
}

// New returns a new instance of a type that fulfills the service.Module
// interface and is capable of provisioning Azure Service Bus
func New(
	armDeployer arm.Deployer,
	namespacesClient servicebusSDK.NamespacesClient,
	queuesClient servicebusSDK.QueuesClient,
	topicsClient servicebusSDK.TopicsClient,
	subscriptionsClient servicebusSDK.SubscriptionsClient,
) service.Module {
	return &module{
		namespaceManager: &namespaceManager{
			armDeployer:      armDeployer,
			namespacesClient: namespacesClient,
		},
		queueManager: &queueManager{
			armDeployer:  armDeployer,
			queuesClient: queuesClient,
		},
		topicManager: &topicManager{
			armDeployer:         armDeployer,
			topicsClient:        topicsClient,
			subscriptionsClient: subscriptionsClient,
		},
	}
}

func (m *module) GetName() string {
	return "servicebus"
}

func (m *module) GetStability() service.Stability {
	return service.StabilityExperimental
}

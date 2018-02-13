package eventhubs

import (
	eventHubSDK "github.com/Azure/azure-sdk-for-go/services/eventhub/mgmt/2017-04-01/eventhub" // nolint: lll
	"github.com/Azure/open-service-broker-azure/pkg/azure/arm"
	"github.com/Azure/open-service-broker-azure/pkg/service"
)

type module struct {
	serviceManager *serviceManager
}

type serviceManager struct {
	armDeployer      arm.Deployer
	namespacesClient eventHubSDK.NamespacesClient
}

// New returns a new instance of a type that fulfills the service.Module
// interface and is capable of provisioning Azure Event Hub
func New(
	armDeployer arm.Deployer,
	namespacesClient eventHubSDK.NamespacesClient,
) service.Module {
	return &module{
		serviceManager: &serviceManager{
			armDeployer:      armDeployer,
			namespacesClient: namespacesClient,
		},
	}
}

func (m *module) GetName() string {
	return "eventhub"
}

func (m *module) GetStability() service.Stability {
	return service.StabilityExperimental
}

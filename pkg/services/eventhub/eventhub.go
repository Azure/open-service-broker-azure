package eventhub

import (
	"github.com/Azure/open-service-broker-azure/pkg/azure/arm"
	"github.com/Azure/open-service-broker-azure/pkg/azure/eventhub"
	"github.com/Azure/open-service-broker-azure/pkg/service"
)

type module struct {
	serviceManager *serviceManager
}

type serviceManager struct {
	armDeployer     arm.Deployer
	eventHubManager eventhub.Manager
}

// New returns a new instance of a type that fulfills the service.Module
// interface and is capable of provisioning Azure Event Hub
func New(
	armDeployer arm.Deployer,
	eventHubManager eventhub.Manager,
) service.Module {
	return &module{
		serviceManager: &serviceManager{
			armDeployer:     armDeployer,
			eventHubManager: eventHubManager,
		},
	}
}

func (m *module) GetName() string {
	return "eventhub"
}

func (m *module) GetStability() service.Stability {
	return service.StabilityAlpha
}

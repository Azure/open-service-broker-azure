package eventhub

import (
	"github.com/Azure/azure-service-broker/pkg/azure/arm"
	"github.com/Azure/azure-service-broker/pkg/azure/eventhub"
	"github.com/Azure/azure-service-broker/pkg/service"
)

type module struct {
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
		armDeployer:     armDeployer,
		eventHubManager: eventHubManager,
	}
}

func (m *module) GetName() string {
	return "eventhub"
}

func (m *module) GetStability() service.Stability {
	return service.StabilityAlpha
}

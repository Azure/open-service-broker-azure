package servicebus

import (
	"github.com/Azure/azure-service-broker/pkg/azure/arm"
	"github.com/Azure/azure-service-broker/pkg/azure/servicebus"
	"github.com/Azure/azure-service-broker/pkg/service"
)

type module struct {
	armDeployer       arm.Deployer
	serviceBusManager servicebus.Manager
}

// New returns a new instance of a type that fulfills the service.Module
// interface and is capable of provisioning Azure Service Bus
func New(
	armDeployer arm.Deployer,
	serviceBusManager servicebus.Manager,
) service.Module {
	return &module{
		armDeployer:       armDeployer,
		serviceBusManager: serviceBusManager,
	}
}

func (m *module) GetName() string {
	return "servicebus"
}

func (m *module) GetStability() service.Stability {
	return service.StabilityAlpha
}

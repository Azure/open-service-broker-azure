package servicebus

import (
	"github.com/Azure/open-service-broker-azure/pkg/azure/arm"
	"github.com/Azure/open-service-broker-azure/pkg/azure/servicebus"
	"github.com/Azure/open-service-broker-azure/pkg/service"
)

type module struct {
	serviceManager *serviceManager
}

type serviceManager struct {
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
		serviceManager: &serviceManager{
			armDeployer:       armDeployer,
			serviceBusManager: serviceBusManager,
		},
	}
}

func (m *module) GetName() string {
	return "servicebus"
}

func (m *module) GetStability() service.Stability {
	return service.StabilityExperimental
}

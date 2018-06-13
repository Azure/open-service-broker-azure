// +build experimental

package servicebus

import (
	servicebusSDK "github.com/Azure/azure-sdk-for-go/services/servicebus/mgmt/2017-04-01/servicebus" // nolint: lll
	"github.com/Azure/open-service-broker-azure/pkg/azure/arm"
	"github.com/Azure/open-service-broker-azure/pkg/service"
)

type module struct {
	serviceManager *serviceManager
}

type serviceManager struct {
	armDeployer      arm.Deployer
	namespacesClient servicebusSDK.NamespacesClient
}

// New returns a new instance of a type that fulfills the service.Module
// interface and is capable of provisioning Azure Service Bus
func New(
	armDeployer arm.Deployer,
	namespacesClient servicebusSDK.NamespacesClient,
) service.Module {
	return &module{
		serviceManager: &serviceManager{
			armDeployer:      armDeployer,
			namespacesClient: namespacesClient,
		},
	}
}

func (m *module) GetName() string {
	return "servicebus"
}

func (m *module) GetStability() service.Stability {
	return service.StabilityExperimental
}

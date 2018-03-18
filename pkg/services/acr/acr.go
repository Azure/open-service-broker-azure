package acr

import (
	"github.com/Azure/open-service-broker-azure/pkg/azure/acr"
	"github.com/Azure/open-service-broker-azure/pkg/azure/arm"
	"github.com/Azure/open-service-broker-azure/pkg/service"
)

type module struct {
	serviceManager *serviceManager
}

type serviceManager struct {
	armDeployer arm.Deployer
	acrManager  acr.Manager
}

// New returns a new instance of a type that fulfills the service.Module
// interface and is capable of provisioning Azure acr
func New(
	armDeployer arm.Deployer,
	acrManager acr.Manager,
) service.Module {
	return &module{
		serviceManager: &serviceManager{
			armDeployer: armDeployer,
			acrManager:  acrManager,
		},
	}
}

func (m *module) GetName() string {
	return "azurecontainerregistry"
}

func (m *module) GetStability() service.Stability {
	return service.StabilityExperimental
}

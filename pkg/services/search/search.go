package search

import (
	"github.com/Azure/open-service-broker-azure/pkg/azure/arm"
	"github.com/Azure/open-service-broker-azure/pkg/azure/search"
	"github.com/Azure/open-service-broker-azure/pkg/service"
)

type module struct {
	serviceManager *serviceManager
}

type serviceManager struct {
	armDeployer   arm.Deployer
	searchManager search.Manager
}

// New returns a new instance of a type that fulfills the service.Module
// interface and is capable of provisioning Azure Search
func New(
	armDeployer arm.Deployer,
	searchManager search.Manager,
) service.Module {
	return &module{
		serviceManager: &serviceManager{
			armDeployer:   armDeployer,
			searchManager: searchManager,
		},
	}
}

func (m *module) GetName() string {
	return "azuresearch"
}

func (m *module) GetStability() service.Stability {
	return service.StabilityAlpha
}

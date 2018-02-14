package aci

import (
	"github.com/Azure/azure-sdk-for-go/arm/containerinstance"
	"github.com/Azure/open-service-broker-azure/pkg/azure/arm"
	"github.com/Azure/open-service-broker-azure/pkg/service"
)

type module struct {
	serviceManager *serviceManager
}

type serviceManager struct {
	armDeployer arm.Deployer
	aciClient   containerinstance.ContainerGroupsClient
}

// New returns a new instance of a type that fulfills the service.Module
// interface and is capable of provisioning Key Vault using "Azure Key Vault"
func New(
	armDeployer arm.Deployer,
	aciClient containerinstance.ContainerGroupsClient,
) service.Module {
	return &module{
		serviceManager: &serviceManager{
			armDeployer: armDeployer,
			aciClient:   aciClient,
		},
	}
}

func (m *module) GetName() string {
	return "aci"
}

func (m *module) GetStability() service.Stability {
	return service.StabilityExperimental
}

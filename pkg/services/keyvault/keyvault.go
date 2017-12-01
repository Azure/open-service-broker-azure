package keyvault

import (
	"github.com/Azure/azure-service-broker/pkg/azure/arm"
	"github.com/Azure/azure-service-broker/pkg/azure/keyvault"
	"github.com/Azure/azure-service-broker/pkg/service"
)

type module struct {
	serviceManager *serviceManager
}

type serviceManager struct {
	armDeployer     arm.Deployer
	keyvaultManager keyvault.Manager
}

// New returns a new instance of a type that fulfills the service.Module
// interface and is capable of provisioning Key Vault using "Azure Key Vault"
func New(
	armDeployer arm.Deployer,
	keyvaultManager keyvault.Manager,
) service.Module {
	return &module{
		serviceManager: &serviceManager{
			armDeployer:     armDeployer,
			keyvaultManager: keyvaultManager,
		},
	}
}

func (m *module) GetName() string {
	return "keyvault"
}

func (m *module) GetStability() service.Stability {
	return service.StabilityAlpha
}

package aci

import (
	"github.com/Azure/open-service-broker-azure/pkg/azure/aci"
	"github.com/Azure/open-service-broker-azure/pkg/azure/arm"
	"github.com/Azure/open-service-broker-azure/pkg/service"
)

type module struct {
	serviceManager *serviceManager
}

type serviceManager struct {
	armDeployer arm.Deployer
	aciManager  aci.Manager
}

// New returns a new instance of a type that fulfills the service.Module
// interface and is capable of provisioning Key Vault using "Azure Key Vault"
func New(
	armDeployer arm.Deployer,
	aciManager aci.Manager,
) service.Module {
	return &module{
		serviceManager: &serviceManager{
			armDeployer: armDeployer,
			aciManager:  aciManager,
		},
	}
}

func (m *module) GetName() string {
	return "aci"
}

func (m *module) GetStability() service.Stability {
	return service.StabilityExperimental
}

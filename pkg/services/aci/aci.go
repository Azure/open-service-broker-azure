package aci

import (
	"github.com/Azure/azure-service-broker/pkg/azure/aci"
	"github.com/Azure/azure-service-broker/pkg/azure/arm"
	"github.com/Azure/azure-service-broker/pkg/service"
)

type module struct {
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
		armDeployer: armDeployer,
		aciManager:  aciManager,
	}
}

func (m *module) GetName() string {
	return "aci"
}

func (m *module) GetStability() service.Stability {
	return service.StabilityAlpha
}

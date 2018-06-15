// +build experimental

package aci

import (
	aciSDK "github.com/Azure/azure-sdk-for-go/services/containerinstance/mgmt/2017-08-01-preview/containerinstance" // nolint: lll
	"open-service-broker-azure/pkg/azure/arm"
	"open-service-broker-azure/pkg/service"
)

type module struct {
	serviceManager *serviceManager
}

type serviceManager struct {
	armDeployer arm.Deployer
	aciClient   aciSDK.ContainerGroupsClient
}

// New returns a new instance of a type that fulfills the service.Module
// interface and is capable of provisioning Key Vault using "Azure Key Vault"
func New(
	armDeployer arm.Deployer,
	aciClient aciSDK.ContainerGroupsClient,
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

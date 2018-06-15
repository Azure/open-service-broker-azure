// +build experimental

package keyvault

import (
	keyVaultSDK "github.com/Azure/azure-sdk-for-go/services/keyvault/mgmt/2016-10-01/keyvault" // nolint: lll
	"open-service-broker-azure/pkg/azure/arm"
	"open-service-broker-azure/pkg/service"
)

type module struct {
	serviceManager *serviceManager
}

type serviceManager struct {
	tenantID     string
	armDeployer  arm.Deployer
	vaultsClient keyVaultSDK.VaultsClient
}

// New returns a new instance of a type that fulfills the service.Module
// interface and is capable of provisioning Key Vault using "Azure Key Vault"
func New(
	tenantID string,
	armDeployer arm.Deployer,
	vaultsClient keyVaultSDK.VaultsClient,
) service.Module {
	return &module{
		serviceManager: &serviceManager{
			tenantID:     tenantID,
			armDeployer:  armDeployer,
			vaultsClient: vaultsClient,
		},
	}
}

func (m *module) GetName() string {
	return "keyvault"
}

func (m *module) GetStability() service.Stability {
	return service.StabilityExperimental
}

package storage

import (
	storageSDK "github.com/Azure/azure-sdk-for-go/arm/storage"
	"github.com/Azure/open-service-broker-azure/pkg/azure/arm"
	"github.com/Azure/open-service-broker-azure/pkg/service"
)

type module struct {
	serviceManager *serviceManager
}

type serviceManager struct {
	armDeployer    arm.Deployer
	accountsClient storageSDK.AccountsClient
}

// New returns a new instance of a type that fulfills the service.Module
// interface and is capable of provisioning Storage using "Azure Storage"
func New(
	armDeployer arm.Deployer,
	accountsClient storageSDK.AccountsClient,
) service.Module {
	return &module{
		serviceManager: &serviceManager{
			armDeployer:    armDeployer,
			accountsClient: accountsClient,
		},
	}
}

func (m *module) GetName() string {
	return "storage"
}

func (m *module) GetStability() service.Stability {
	return service.StabilityExperimental
}

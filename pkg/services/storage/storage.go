package storage

import (
	"github.com/Azure/open-service-broker-azure/pkg/azure/arm"
	"github.com/Azure/open-service-broker-azure/pkg/azure/storage"
	"github.com/Azure/open-service-broker-azure/pkg/service"
)

type module struct {
	serviceManager *serviceManager
}

type serviceManager struct {
	armDeployer    arm.Deployer
	storageManager storage.Manager
}

// New returns a new instance of a type that fulfills the service.Module
// interface and is capable of provisioning Storage using "Azure Storage"
func New(
	armDeployer arm.Deployer,
	storageManager storage.Manager,
) service.Module {
	return &module{
		serviceManager: &serviceManager{
			armDeployer:    armDeployer,
			storageManager: storageManager,
		},
	}
}

func (m *module) GetName() string {
	return "storage"
}

func (m *module) GetStability() service.Stability {
	return service.StabilityExperimental
}

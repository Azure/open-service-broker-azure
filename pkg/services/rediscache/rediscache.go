package rediscache

import (
	redisSDK "github.com/Azure/azure-sdk-for-go/arm/redis"
	"github.com/Azure/open-service-broker-azure/pkg/azure/arm"
	"github.com/Azure/open-service-broker-azure/pkg/service"
)

type module struct {
	serviceManager *serviceManager
}

type serviceManager struct {
	armDeployer arm.Deployer
	groupClient redisSDK.GroupClient
}

// New returns a new instance of a type that fulfills the service.Module
// interface and is capable of provisioning Redis using "Azure Redis Cache"
func New(
	armDeployer arm.Deployer,
	groupClient redisSDK.GroupClient,
) service.Module {
	return &module{
		serviceManager: &serviceManager{
			armDeployer: armDeployer,
			groupClient: groupClient,
		},
	}
}

func (m *module) GetName() string {
	return "rediscache"
}

func (m *module) GetStability() service.Stability {
	return service.StabilityExperimental
}

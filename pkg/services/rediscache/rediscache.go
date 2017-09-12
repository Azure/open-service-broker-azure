package rediscache

import (
	"github.com/Azure/azure-service-broker/pkg/azure/arm"
	"github.com/Azure/azure-service-broker/pkg/azure/rediscache"
	"github.com/Azure/azure-service-broker/pkg/service"
)

type module struct {
	armDeployer  arm.Deployer
	redisManager rediscache.Manager
}

// New returns a new instance of a type that fulfills the service.Module
// interface and is capable of provisioning Redis using "Azure Redis Cache"
func New(
	armDeployer arm.Deployer,
	redisManager rediscache.Manager,
) service.Module {
	return &module{
		armDeployer:  armDeployer,
		redisManager: redisManager,
	}
}

func (m *module) GetName() string {
	return "rediscache"
}

func (m *module) GetStability() service.Stability {
	return service.StabilityAlpha
}

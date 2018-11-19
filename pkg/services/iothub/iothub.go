package iothub

import (
	iotHubSDK "github.com/Azure/azure-sdk-for-go/services/iothub/mgmt/2017-07-01/devices"
	"github.com/Azure/open-service-broker-azure/pkg/azure/arm"
	"github.com/Azure/open-service-broker-azure/pkg/service"
)

type module struct {
	iotHubManager *iotHubManager
}

type iotHubManager struct {
	armDeployer  arm.Deployer
	iotHubClient iotHubSDK.IotHubResourceClient
}

// New returns a new instance of a type that fulfills the service.Module
// interface and is capable of provisioning text analytics using
// "Azure Cognitive Services"
func New(
	armDeployer arm.Deployer,
	iotHubClient iotHubSDK.IotHubResourceClient,
) service.Module {
	return &module{
		iotHubManager: &iotHubManager{
			armDeployer:  armDeployer,
			iotHubClient: iotHubClient,
		},
	}
}

func (m *module) GetName() string {
	return "iotHub"
}

func (m *module) GetStability() service.Stability {
	return service.StabilityExperimental
}

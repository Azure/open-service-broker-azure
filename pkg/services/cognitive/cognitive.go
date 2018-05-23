package cognitive

import (
	cognitiveSDK "github.com/Azure/azure-sdk-for-go/services/cognitiveservices/mgmt/2017-04-18/cognitiveservices" // nolint: lll
	"github.com/Azure/open-service-broker-azure/pkg/azure/arm"
	"github.com/Azure/open-service-broker-azure/pkg/service"
)

type module struct {
	serviceManager *serviceManager
}

type serviceManager struct {
	armDeployer      arm.Deployer
	congnitiveClient cognitiveSDK.AccountsClient
}

// New returns a new instance of a type that fulfills the service.Module
// interface and is capable of provisioning text analytics using
// "Azure Cognitive Services"
func New(
	armDeployer arm.Deployer,
	congnitiveClient cognitiveSDK.AccountsClient,
) service.Module {
	return &module{
		serviceManager: &serviceManager{
			armDeployer:      armDeployer,
			congnitiveClient: congnitiveClient,
		},
	}
}

func (m *module) GetName() string {
	return "cognitive"
}

func (m *module) GetStability() service.Stability {
	return service.StabilityExperimental
}

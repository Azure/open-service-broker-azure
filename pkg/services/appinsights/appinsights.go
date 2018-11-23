package appinsights

import (
	appInsightsSDK "github.com/Azure/azure-sdk-for-go/services/appinsights/mgmt/2015-05-01/insights" // nolint: lll
	"github.com/Azure/open-service-broker-azure/pkg/azure/arm"
	"github.com/Azure/open-service-broker-azure/pkg/service"
)

type module struct {
	serviceManager *serviceManager
}

type serviceManager struct {
	armDeployer       arm.Deployer
	appInsightsClient appInsightsSDK.ComponentsClient
}

// New returns a new instance of a type that fulfills the service.Module
// interface and is capable of provisioning Azure Application Insights
func New(
	armDeployer arm.Deployer,
	appInsightsClient appInsightsSDK.ComponentsClient,
) service.Module {
	return &module{
		serviceManager: &serviceManager{
			armDeployer:       armDeployer,
			appInsightsClient: appInsightsClient,
		},
	}
}

func (m *module) GetName() string {
	return "appinsights"
}

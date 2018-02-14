package postgresqldb

import (
	postgresSDK "github.com/Azure/azure-sdk-for-go/services/postgresql/mgmt/2017-04-30-preview/postgresql" // nolint: lll
	"github.com/Azure/open-service-broker-azure/pkg/azure/arm"
	"github.com/Azure/open-service-broker-azure/pkg/service"
)

type module struct {
	serviceManager *serviceManager
}

type serviceManager struct {
	armDeployer                 arm.Deployer
	checkNameAvailabilityClient postgresSDK.CheckNameAvailabilityClient
	serversClient               postgresSDK.ServersClient
}

// New returns a new instance of a type that fulfills the service.Module
// interface and is capable of provisioning PostgreSQL servers and databases
// using "Azure Database for PostgreSQL"
func New(
	armDeployer arm.Deployer,
	checkNameAvailabilityClient postgresSDK.CheckNameAvailabilityClient,
	serversClient postgresSDK.ServersClient,
) service.Module {
	return &module{
		serviceManager: &serviceManager{
			armDeployer:                 armDeployer,
			checkNameAvailabilityClient: checkNameAvailabilityClient,
			serversClient:               serversClient,
		},
	}
}

func (m *module) GetName() string {
	return "postgresql"
}

func (m *module) GetStability() service.Stability {
	return service.StabilityExperimental
}

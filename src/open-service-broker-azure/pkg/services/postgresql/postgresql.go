package postgresql

import (
	postgresSDK "github.com/Azure/azure-sdk-for-go/services/postgresql/mgmt/2017-04-30-preview/postgresql" // nolint: lll
	"open-service-broker-azure/pkg/azure/arm"
	"open-service-broker-azure/pkg/service"
)

type module struct {
	allInOneManager *allInOneManager
	databaseManager *databaseManager
	dbmsManager     *dbmsManager
}

type allInOneManager struct {
	armDeployer                 arm.Deployer
	checkNameAvailabilityClient postgresSDK.CheckNameAvailabilityClient
	serversClient               postgresSDK.ServersClient
}

type databaseManager struct {
	armDeployer                 arm.Deployer
	checkNameAvailabilityClient postgresSDK.CheckNameAvailabilityClient
	databasesClient             postgresSDK.DatabasesClient
}

type dbmsManager struct {
	armDeployer                 arm.Deployer
	checkNameAvailabilityClient postgresSDK.CheckNameAvailabilityClient
	serversClient               postgresSDK.ServersClient
}

// New returns a new instance of a type that fulfills the service.Module
// interface and is capable of provisioning PostgreSQL DBMS and databases
// using "Azure Database for PostgreSQL"
func New(
	armDeployer arm.Deployer,
	checkNameAvailabilityClient postgresSDK.CheckNameAvailabilityClient,
	serversClient postgresSDK.ServersClient,
	databasesClient postgresSDK.DatabasesClient,
) service.Module {
	return &module{
		allInOneManager: &allInOneManager{
			armDeployer:                 armDeployer,
			checkNameAvailabilityClient: checkNameAvailabilityClient,
			serversClient:               serversClient,
		},
		databaseManager: &databaseManager{
			armDeployer:                 armDeployer,
			checkNameAvailabilityClient: checkNameAvailabilityClient,
			databasesClient:             databasesClient,
		},
		dbmsManager: &dbmsManager{
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
	return service.StabilityPreview
}

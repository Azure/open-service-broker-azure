package mssql

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/arm/sql"
	az "github.com/Azure/open-service-broker-azure/pkg/azure"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Manager is an interface to be implemented by any component capable of
// managing Azure SQL Database
type Manager interface {
	DeleteServer(
		serverName string,
		resourceGroupName string,
	) error
	DeleteDatabase(
		serverName string,
		databaseName string,
		resourceGroupName string,
	) error
}

type manager struct {
	azureEnvironment azure.Environment
	subscriptionID   string
	tenantID         string
	clientID         string
	clientSecret     string
}

// NewManager returns a new implementation of the Manager interface
func NewManager() (Manager, error) {
	azureConfig, err := az.GetConfig()
	if err != nil {
		return nil, err
	}
	azureEnvironment, err := azure.EnvironmentFromName(azureConfig.Environment)
	if err != nil {
		return nil, fmt.Errorf(
			`error parsing Azure environment name "%s"`,
			azureConfig.Environment,
		)
	}
	return &manager{
		azureEnvironment: azureEnvironment,
		subscriptionID:   azureConfig.SubscriptionID,
		tenantID:         azureConfig.TenantID,
		clientID:         azureConfig.ClientID,
		clientSecret:     azureConfig.ClientSecret,
	}, nil
}

func (m *manager) DeleteServer(
	serverName string,
	resourceGroupName string,
) error {
	authorizer, err := az.GetBearerTokenAuthorizer(
		m.azureEnvironment,
		m.tenantID,
		m.clientID,
		m.clientSecret,
	)
	if err != nil {
		return fmt.Errorf("error getting bearer token authorizer: %s", err)
	}

	serversClient := sql.NewServersClientWithBaseURI(
		m.azureEnvironment.ResourceManagerEndpoint,
		m.subscriptionID,
	)
	serversClient.Authorizer = authorizer
	if _, err = serversClient.Delete(
		resourceGroupName,
		serverName,
	); err != nil {
		return fmt.Errorf("error deleting mssql server: %s", err)
	}

	return nil
}

func (m *manager) DeleteDatabase(
	serverName string,
	databaseName string,
	resourceGroupName string,
) error {
	authorizer, err := az.GetBearerTokenAuthorizer(
		m.azureEnvironment,
		m.tenantID,
		m.clientID,
		m.clientSecret,
	)
	if err != nil {
		return fmt.Errorf("error getting bearer token authorizer: %s", err)
	}

	databasesClient := sql.NewDatabasesClientWithBaseURI(
		m.azureEnvironment.ResourceManagerEndpoint,
		m.subscriptionID,
	)
	databasesClient.Authorizer = authorizer
	if _, err = databasesClient.Delete(
		resourceGroupName,
		serverName,
		databaseName,
	); err != nil {
		return fmt.Errorf("error deleting mssql database: %s", err)
	}

	return nil
}

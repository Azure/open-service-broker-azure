package rediscache

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/arm/redis"
	az "github.com/Azure/azure-service-broker/pkg/azure"
	"github.com/Azure/go-autorest/autorest/azure"
)

// krancour: This whole interface and its implementations are a workaround!
// Ideally, we'd like to accomplish as much as possible using ARM, however,
// currently, when deleting ARM deployments, Redis don't get
// deleted. This seems to be an issue with the underlying RP. To work around
// that, we'll use this interface and implementations thereof to delete those
// servers directly.

// Manager is an interface to be implemented by any component capable of
// managing Azure Redis Cache
type Manager interface {
	DeleteServer(
		serverName string,
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

	serversClient := redis.NewGroupClientWithBaseURI(
		m.azureEnvironment.ResourceManagerEndpoint,
		m.subscriptionID,
	)
	serversClient.Authorizer = authorizer
	cancelCh := make(chan struct{})
	_, errChan := serversClient.Delete(
		resourceGroupName,
		serverName,
		cancelCh,
	)
	if err := <-errChan; err != nil {
		return fmt.Errorf("error deleting redis server: %s", err)
	}

	return nil
}

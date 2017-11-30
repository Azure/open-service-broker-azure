package servicebus

import (
	"fmt"
	"strings"

	"github.com/Azure/azure-sdk-for-go/arm/servicebus"
	az "github.com/Azure/open-service-broker-azure/pkg/azure"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Manager is an interface to be implemented by any component capable of
// managing Azure Service Bus
type Manager interface {
	DeleteNamespace(
		serviceBusNamespaceName string,
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

func (m *manager) DeleteNamespace(
	serviceBusNamespaceName string,
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

	nsClient := servicebus.NewNamespacesClientWithBaseURI(
		m.azureEnvironment.ResourceManagerEndpoint,
		m.subscriptionID,
	)
	nsClient.Authorizer = authorizer
	cancelCh := make(chan struct{})
	_, errChan := nsClient.Delete(
		resourceGroupName,
		serviceBusNamespaceName,
		cancelCh,
	)
	if err := <-errChan; err != nil {
		// Workaround for https://github.com/Azure/azure-sdk-for-go/issues/759
		if strings.Contains(err.Error(), "StatusCode=404") {
			return nil
		}
		return fmt.Errorf("error deleting azure service bus namespace: %s", err)
	}

	return nil
}

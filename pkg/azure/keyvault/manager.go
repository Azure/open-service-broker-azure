package keyvault

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/arm/keyvault"
	az "github.com/Azure/open-service-broker-azure/pkg/azure"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Manager is an interface to be implemented by any component capable of
// managing Azure Key Vault
type Manager interface {
	GetTenantID() string

	DeleteVault(
		vaultName string,
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

func (m *manager) GetTenantID() string {
	return m.tenantID
}

func (m *manager) DeleteVault(
	vaultName string,
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

	vaultClient := keyvault.NewVaultsClientWithBaseURI(
		m.azureEnvironment.ResourceManagerEndpoint,
		m.subscriptionID,
	)
	vaultClient.Authorizer = authorizer
	_, err = vaultClient.Delete(
		resourceGroupName,
		vaultName,
	)
	if err != nil {
		return fmt.Errorf("error deleting keyvault: %s", err)
	}

	return nil
}

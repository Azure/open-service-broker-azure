package hdinsight

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/arm/hdinsight"
	"github.com/Azure/azure-sdk-for-go/arm/storage"
	"github.com/Azure/azure-service-broker/pkg/azure"
	az "github.com/Azure/go-autorest/autorest/azure"
)

// Manager is an interface to be implemented by any component capable of
// managing Azure HDInsight and the additional Azure Storage
type Manager interface {
	DeleteCluster(
		clusterName string,
		resourceGroupName string,
	) error
	DeleteStorageAccount(
		storageAccountName string,
		resourceGroupName string,
	) error
}

type manager struct {
	azureEnvironment az.Environment
	subscriptionID   string
	tenantID         string
	clientID         string
	clientSecret     string
}

// NewManager returns a new implementation of the Manager interface
func NewManager() (Manager, error) {
	azureConfig, err := azure.GetConfig()
	if err != nil {
		return nil, err
	}
	azureEnvironment, err := az.EnvironmentFromName(azureConfig.Environment)
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

func (m *manager) DeleteCluster(
	clusterName string,
	resourceGroupName string,
) error {
	authorizer, err := azure.GetBearerTokenAuthorizer(
		m.azureEnvironment,
		m.tenantID,
		m.clientID,
		m.clientSecret,
	)
	if err != nil {
		return fmt.Errorf("error getting bearer token authorizer: %s", err)
	}

	clustersClient := hdinsight.NewClustersClientWithBaseURI(
		m.azureEnvironment.ResourceManagerEndpoint,
		m.subscriptionID,
	)
	clustersClient.Authorizer = authorizer
	cancelCh := make(chan struct{})
	_, errChan := clustersClient.Delete(
		resourceGroupName,
		clusterName,
		cancelCh,
	)
	if err := <-errChan; err != nil {
		return fmt.Errorf("error deleting hdinsight cluster: %s", err)
	}

	return nil
}

func (m *manager) DeleteStorageAccount(
	storageAccountName string,
	resourceGroupName string,
) error {
	authorizer, err := azure.GetBearerTokenAuthorizer(
		m.azureEnvironment,
		m.tenantID,
		m.clientID,
		m.clientSecret,
	)
	if err != nil {
		return fmt.Errorf("error getting bearer token authorizer: %s", err)
	}

	accountsClient := storage.NewAccountsClientWithBaseURI(
		m.azureEnvironment.ResourceManagerEndpoint,
		m.subscriptionID,
	)
	accountsClient.Authorizer = authorizer

	if _, err := accountsClient.Delete(
		resourceGroupName,
		storageAccountName,
	); err != nil {
		return fmt.Errorf("error deleting hdinsight storage account: %s", err)
	}

	return nil
}

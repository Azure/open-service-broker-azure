package hdinsight

import (
	"fmt"

	"github.com/Azure/azure-service-broker/pkg/azure"
	"github.com/Azure/azure-service-broker/pkg/service"
)

func (m *module) ValidateBindingParameters(
	bindingParameters service.BindingParameters,
) error {
	return nil
}

func (m *module) Bind(
	provisioningContext service.ProvisioningContext,
	bindingParameters service.BindingParameters,
) (service.BindingContext, service.Credentials, error) {
	pc, ok := provisioningContext.(*hdinsightProvisioningContext)
	if !ok {
		return nil, nil, fmt.Errorf(
			"error casting provisioningContext as hdinsightProvisioningContext",
		)
	}

	var clusterEndpoint string

	azureConfig, err := azure.GetConfig()
	if err != nil {
		return nil, nil, err
	}
	switch azureConfig.Environment {
	case "AzurePublicCloud":
		clusterEndpoint = pc.ClusterName + ".azurehdinsight.net"
	case "AzureUSGovernmentCloud":
		clusterEndpoint = pc.ClusterName + ".azurehdinsight.us"
	case "AzureChinaCloud":
		clusterEndpoint = pc.ClusterName + ".azurehdinsight.cn"
	case "AzureGermanCloud":
		clusterEndpoint = pc.ClusterName + ".azurehdinsight.de"
	default:
		return nil, nil, fmt.Errorf(
			"error unknown cluster endpoint for the environment",
		)
	}
	return &hdinsightBindingContext{},
		&hdinsightCredentials{
			ClusterEndpoint:          clusterEndpoint,
			Username:                 pc.ClusterLoginUserName,
			Password:                 pc.ClusterLoginPassword,
			StorageAccountName:       pc.StorageAccountName,
			StorageAccountKey:        pc.StorageAccountKey,
			BlobStorageEndpoint:      pc.BlobStorageEndpoint,
			BlobStorageContainerName: pc.BlobStorageContainerName,
		},
		nil
}

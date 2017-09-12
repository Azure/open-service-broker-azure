package hdinsight

import "github.com/Azure/azure-service-broker/pkg/service"

// ProvisioningParameters encapsulates HDInsight-specific provisioning options
type ProvisioningParameters struct {
	Location               string            `json:"location"`
	ResourceGroup          string            `json:"resourceGroup"`
	ClusterWorkerNodeCount int               `json:"clusterWorkerNodeCount"`
	Tags                   map[string]string `json:"tags"`
}

type hdinsightProvisioningContext struct {
	ResourceGroupName        string `json:"resourceGroup"`
	ARMDeploymentName        string `json:"armDeployment"`
	ClusterName              string `json:"clusterName"`
	ClusterLoginUserName     string `json:"clusterLoginUserName"`
	ClusterLoginPassword     string `json:"clusterLoginPassword"`
	SSHUserName              string `json:"sshUserName"`
	SSHPassword              string `json:"sshPassword"`
	StorageAccountName       string `json:"storageAccountName"`
	StorageAccountKey        string `json:"storageAccountKey"`
	BlobStorageEndpoint      string `json:"blobStorageEndpoint"`
	BlobStorageContainerName string `json:"blobStorageContainerName"`
}

// BindingParameters encapsulates HDInsight-specific binding options
type BindingParameters struct {
}

type hdinsightBindingContext struct {
}

type hdinsightCredentials struct {
	ClusterEndpoint          string `json:"clusterEndpoint"`
	Username                 string `json:"username"`
	Password                 string `json:"password"`
	StorageAccountName       string `json:"storageAccountName"`
	StorageAccountKey        string `json:"storageAccountKey"`
	BlobStorageEndpoint      string `json:"blobStorageEndpoint"`
	BlobStorageContainerName string `json:"blobStorageContainerName"`
}

func (m *hdinsightProvisioningContext) GetResourceGroupName() string {
	return m.ResourceGroupName
}

func (
	m *module,
) GetEmptyProvisioningParameters() service.ProvisioningParameters {
	return &ProvisioningParameters{}
}

func (m *module) GetEmptyProvisioningContext() service.ProvisioningContext {
	return &hdinsightProvisioningContext{}
}

func (m *module) GetEmptyBindingParameters() service.BindingParameters {
	return &BindingParameters{}
}

func (m *module) GetEmptyBindingContext() service.BindingContext {
	return &hdinsightBindingContext{}
}

func (m *module) GetEmptyCredentials() service.Credentials {
	return &hdinsightCredentials{}
}

package hdinsight

import "github.com/Azure/open-service-broker-azure/pkg/service"

// ProvisioningParameters encapsulates HDInsight-specific provisioning options
type ProvisioningParameters struct {
	ClusterWorkerNodeCount int               `json:"clusterWorkerNodeCount"`
}

type hdinsightProvisioningContext struct {
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

// UpdatingParameters encapsulates HDInsight-specific updating options
type UpdatingParameters struct {
}

// BindingParameters encapsulates HDInsight-specific binding options
type BindingParameters struct {
}

type hdinsightBindingContext struct {
}

// Credentials encapsulates HDInsight-specific coonection details and
// credentials.
type Credentials struct {
	ClusterEndpoint          string `json:"clusterEndpoint"`
	Username                 string `json:"username"`
	Password                 string `json:"password"`
	StorageAccountName       string `json:"storageAccountName"`
	StorageAccountKey        string `json:"storageAccountKey"`
	BlobStorageEndpoint      string `json:"blobStorageEndpoint"`
	BlobStorageContainerName string `json:"blobStorageContainerName"`
}

func (
	s *serviceManager,
) GetEmptyProvisioningParameters() service.ProvisioningParameters {
	return &ProvisioningParameters{}
}

func (
	s *serviceManager,
) GetEmptyProvisioningContext() service.ProvisioningContext {
	return &hdinsightProvisioningContext{}
}

func (s *serviceManager) GetEmptyBindingParameters() service.BindingParameters {
	return &BindingParameters{}
}

func (s *serviceManager) GetEmptyBindingContext() service.BindingContext {
	return &hdinsightBindingContext{}
}

func (s *serviceManager) GetEmptyCredentials() service.Credentials {
	return &hdinsightCredentials{}
}

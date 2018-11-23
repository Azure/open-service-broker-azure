package storage

import "github.com/Azure/open-service-broker-azure/pkg/service"

type instanceDetails struct {
	ARMDeploymentName  string `json:"armDeployment"`
	StorageAccountName string `json:"storageAccountName"`
	ContainerName      string `json:"containerName"`
	AccessKey          string `json:"accessKey"`
}

type credentials struct {
	StorageAccountName          string `json:"storageAccountName"`
	AccessKey                   string `json:"accessKey"`
	PrimaryBlobServiceEndPoint  string `json:"primaryBlobServiceEndPoint,omitempty"`  // nolint: lll
	PrimaryFileServiceEndPoint  string `json:"primaryFileServiceEndPoint,omitempty"`  // nolint: lll
	PrimaryQueueServiceEndPoint string `json:"primaryQueueServiceEndPoint,omitempty"` // nolint: lll
	PrimaryTableServiceEndPoint string `json:"primaryTableServiceEndPoint,omitempty"` // nolint: lll
	ContainerName               string `json:"containerName,omitempty"`
}

func (s *storageManager) GetEmptyInstanceDetails() service.InstanceDetails {
	return &instanceDetails{}
}

func (s *storageManager) GetEmptyBindingDetails() service.BindingDetails {
	return nil
}

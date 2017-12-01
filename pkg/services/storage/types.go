package storage

import "github.com/Azure/azure-service-broker/pkg/service"

type storageKind string

const (
	storageKindGeneralPurposeStorageAcccount storageKind = "GeneralPurposeStorageAccount" // nolint: lll
	storageKindBlobStorageAccount            storageKind = "BlobStorageAccount"
	storageKindBlobContainer                 storageKind = "BlobContainer"
)

// ProvisioningParameters encapsulates Storage-specific provisioning options
type ProvisioningParameters struct{}

type storageProvisioningContext struct {
	ARMDeploymentName  string `json:"armDeployment"`
	StorageAccountName string `json:"storageAccountName"`
	AccessKey          string `json:"accessKey"`
	ContainerName      string `json:"containerName"`
}

// UpdatingParameters encapsulates Storage-specific updating options
type UpdatingParameters struct {
}

// BindingParameters encapsulates Storage-specific binding options
type BindingParameters struct {
}

type storageBindingContext struct {
}

// Credentials encapsulates Storage-specific coonection details and credentials.
type Credentials struct {
	StorageAccountName string `json:"storageAccountName"`
	AccessKey          string `json:"accessKey"`
	ContainerName      string `json:"containerName,omitempty"`
}

func (
	s *serviceManager,
) GetEmptyProvisioningParameters() service.ProvisioningParameters {
	return &ProvisioningParameters{}
}

func (
	s *serviceManager,
) GetEmptyUpdatingParameters() service.UpdatingParameters {
	return &UpdatingParameters{}
}

func (
	s *serviceManager,
) GetEmptyProvisioningContext() service.ProvisioningContext {
	return &storageProvisioningContext{}
}

func (s *serviceManager) GetEmptyBindingParameters() service.BindingParameters {
	return &BindingParameters{}
}

func (s *serviceManager) GetEmptyBindingContext() service.BindingContext {
	return &storageBindingContext{}
}

func (s *serviceManager) GetEmptyCredentials() service.Credentials {
	return &Credentials{}
}

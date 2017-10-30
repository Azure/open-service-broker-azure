package storage

import "github.com/Azure/azure-service-broker/pkg/service"

type storageKind string

const (
	storageKindGeneralPurposeStorageAcccount storageKind = "GeneralPurposeStorageAccount" // nolint: lll
	storageKindBlobStorageAccount            storageKind = "BlobStorageAccount"
	storageKindBlobContainer                 storageKind = "BlobContainer"
)

// ProvisioningParameters encapsulates Storage-specific provisioning options
type ProvisioningParameters struct {
	Location      string            `json:"location"`
	ResourceGroup string            `json:"resourceGroup"`
	Tags          map[string]string `json:"tags"`
}

type storageProvisioningContext struct {
	ResourceGroupName  string `json:"resourceGroup"`
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
	m *module,
) GetEmptyProvisioningParameters() service.ProvisioningParameters {
	return &ProvisioningParameters{}
}

// SetResourceGroup sets the name of the resource group into which service
// instances will be deployed
func (p *ProvisioningParameters) SetResourceGroup(resourceGroup string) {
	p.ResourceGroup = resourceGroup
}

func (
	m *module,
) GetEmptyUpdatingParameters() service.UpdatingParameters {
	return &UpdatingParameters{}
}

func (m *module) GetEmptyProvisioningContext() service.ProvisioningContext {
	return &storageProvisioningContext{}
}

func (m *module) GetEmptyBindingParameters() service.BindingParameters {
	return &BindingParameters{}
}

func (m *module) GetEmptyBindingContext() service.BindingContext {
	return &storageBindingContext{}
}

func (m *module) GetEmptyCredentials() service.Credentials {
	return &Credentials{}
}

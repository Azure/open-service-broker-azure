package cosmosdb

import "github.com/Azure/open-service-broker-azure/pkg/service"

// CosmosCredentials encapsulates CosmosDB-specific for connecting via
// a variety of APIs. This excludes MongoDB.
type CosmosCredentials struct {
	URI                     string `json:"uri,omitempty"`
	PrimaryConnectionString string `json:"primaryConnectionString,omitempty"`
	PrimaryKey              string `json:"primaryKey,omitempty"`
}

func (
	c *cosmosManager,
) GetEmptyProvisioningParameters() service.ProvisioningParameters {
	return &ProvisioningParameters{}
}

func (
	c *cosmosManager,
) GetEmptySecureProvisioningParameters() service.SecureProvisioningParameters {
	return &SecureProvisioningParameters{}
}

func (
	c *cosmosManager,
) GetEmptyInstanceDetails() service.InstanceDetails {
	return &cosmosdbInstanceDetails{}
}

func (
	c *cosmosManager,
) GetEmptySecureInstanceDetails() service.SecureInstanceDetails {
	return &cosmosdbSecureInstanceDetails{}
}

func (c *cosmosManager) GetEmptyBindingParameters() service.BindingParameters {
	return &BindingParameters{}
}

func (
	c *cosmosManager,
) GetEmptySecureBindingParameters() service.SecureBindingParameters {
	return &SecureBindingParameters{}
}

func (c *cosmosManager) GetEmptyBindingDetails() service.BindingDetails {
	return &cosmosdbBindingDetails{}
}

func (
	c *cosmosManager,
) GetEmptySecureBindingDetails() service.SecureBindingDetails {
	return &cosmosdbSecureBindingDetails{}
}

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
	c *cosmosAccountManager,
) GetEmptyProvisioningParameters() service.ProvisioningParameters {
	return nil
}

func (
	c *cosmosAccountManager,
) GetEmptySecureProvisioningParameters() service.SecureProvisioningParameters {
	return nil
}

func (
	c *cosmosAccountManager,
) GetEmptyInstanceDetails() service.InstanceDetails {
	return &cosmosdbInstanceDetails{}
}

func (
	c *cosmosAccountManager,
) GetEmptySecureInstanceDetails() service.SecureInstanceDetails {
	return &cosmosdbSecureInstanceDetails{}
}

func (
	c *cosmosAccountManager,
) GetEmptyBindingParameters() service.BindingParameters {
	return nil
}

func (
	c *cosmosAccountManager,
) GetEmptySecureBindingParameters() service.SecureBindingParameters {
	return nil
}

func (c *cosmosAccountManager) GetEmptyBindingDetails() service.BindingDetails {
	return nil
}

func (
	c *cosmosAccountManager,
) GetEmptySecureBindingDetails() service.SecureBindingDetails {
	return nil
}

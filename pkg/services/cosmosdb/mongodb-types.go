package cosmosdb

import "github.com/Azure/open-service-broker-azure/pkg/service"

// MongoCredentials encapsulates CosmosDB-specific connection details and
// credentials for connecting with the MongoDB API.
type MongoCredentials struct {
	Host             string `json:"host,omitempty"`
	Port             int    `json:"port,omitempty"`
	Username         string `json:"username,omitempty"`
	Password         string `json:"password,omitempty"`
	ConnectionString string `json:"connectionString,omitempty"`
}

func (
	m *mongoManager,
) GetEmptyProvisioningParameters() service.ProvisioningParameters {
	return &ProvisioningParameters{}
}

func (
	m *mongoManager,
) GetEmptySecureProvisioningParameters() service.SecureProvisioningParameters {
	return &SecureProvisioningParameters{}
}

func (
	m *mongoManager,
) GetEmptyInstanceDetails() service.InstanceDetails {
	return &cosmosdbInstanceDetails{}
}

func (
	m *mongoManager,
) GetEmptySecureInstanceDetails() service.SecureInstanceDetails {
	return &cosmosdbSecureInstanceDetails{}
}

func (m *mongoManager) GetEmptyBindingParameters() service.BindingParameters {
	return &BindingParameters{}
}

func (
	m *mongoManager,
) GetEmptySecureBindingParameters() service.SecureBindingParameters {
	return &SecureBindingParameters{}
}

func (m *mongoManager) GetEmptyBindingDetails() service.BindingDetails {
	return &cosmosdbBindingDetails{}
}

func (
	m *mongoManager,
) GetEmptySecureBindingDetails() service.SecureBindingDetails {
	return &cosmosdbSecureBindingDetails{}
}

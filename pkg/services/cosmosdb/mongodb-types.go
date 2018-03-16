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
	m *mongoAccountManager,
) GetEmptyProvisioningParameters() service.ProvisioningParameters {
	return nil
}

func (
	m *mongoAccountManager,
) GetEmptySecureProvisioningParameters() service.SecureProvisioningParameters {
	return nil
}

func (
	m *mongoAccountManager,
) GetEmptyInstanceDetails() service.InstanceDetails {
	return &cosmosdbInstanceDetails{}
}

func (
	m *mongoAccountManager,
) GetEmptySecureInstanceDetails() service.SecureInstanceDetails {
	return &cosmosdbSecureInstanceDetails{}
}

func (
	m *mongoAccountManager,
) GetEmptyBindingParameters() service.BindingParameters {
	return nil
}

func (
	m *mongoAccountManager,
) GetEmptySecureBindingParameters() service.SecureBindingParameters {
	return nil
}

func (m *mongoAccountManager) GetEmptyBindingDetails() service.BindingDetails {
	return nil
}

func (
	m *mongoAccountManager,
) GetEmptySecureBindingDetails() service.SecureBindingDetails {
	return nil
}

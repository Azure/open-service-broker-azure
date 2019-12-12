package cosmosdb

import "github.com/Azure/open-service-broker-azure/pkg/service"

// mongoCredentials encapsulates CosmosDB-specific connection details and
// credentials for connecting with the MongoDB API.
type mongoCredentials struct {
	Host             string `json:"host"`
	Port             int    `json:"port"`
	Username         string `json:"username"`
	Password         string `json:"password"`
	ConnectionString string `json:"connectionString"`
	URI              string `json:"uri"`
}

// GetEmptyInstanceDetails returns an "empty" service-specific object that
// can be populated with data during unmarshaling of JSON to an Instance
func (
	m *mongoAccountManager,
) GetEmptyInstanceDetails() service.InstanceDetails {
	return &cosmosdbInstanceDetails{}
}

// GetEmptyBindingDetails returns an "empty" service-specific object that
// can be populated with data during unmarshaling of JSON to a Binding
func (m *mongoAccountManager) GetEmptyBindingDetails() service.BindingDetails {
	return nil
}

const (
	mongoDBVersion36 = "3.6"
	mongoDBVersion32 = "3.2"
)

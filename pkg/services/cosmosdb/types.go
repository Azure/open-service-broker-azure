package cosmosdb

import "github.com/Azure/azure-service-broker/pkg/service"

type databaseKind string

const (
	databaseKindMongoDB          databaseKind = "MongoDB"
	databaseKindGlobalDocumentDB databaseKind = "GlobalDocumentDB"
)

// ProvisioningParameters encapsulates CosmosDB-specific provisioning options
type ProvisioningParameters struct{}

type cosmosdbProvisioningContext struct {
	ARMDeploymentName        string       `json:"armDeployment"`
	DatabaseAccountName      string       `json:"name"`
	DatabaseKind             databaseKind `json:"kind"`
	FullyQualifiedDomainName string       `json:"fullyQualifiedDomainName"`
	ConnectionString         string       `json:"connectionString"`
	PrimaryKey               string       `json:"primaryKey"`
}

// UpdatingParameters encapsulates CosmosDB-specific updating options
type UpdatingParameters struct {
}

// BindingParameters encapsulates CosmosDB-specific binding options
type BindingParameters struct {
}

type cosmosdbBindingContext struct {
}

// Credentials encapsulates CosmosDB-specific connection details and
// credentials. The attributes of this type cover all the attributes possibly
// used by either of CosmosDBs two connections types-- MongoDB or DocumentDB.
type Credentials struct {
	Host                    string `json:"host,omitempty"`
	Port                    int    `json:"port,omitempty"`
	Username                string `json:"username,omitempty"`
	Password                string `json:"password,omitempty"`
	ConnectionString        string `json:"connectionString,omitempty"`
	URI                     string `json:"uri,omitempty"`
	PrimaryConnectionString string `json:"primaryConnectionString,omitempty"`
	PrimaryKey              string `json:"primaryKey,omitempty"`
}

func (
	m *module,
) GetEmptyProvisioningParameters() service.ProvisioningParameters {
	return &ProvisioningParameters{}
}

func (m *module) GetEmptyProvisioningContext() service.ProvisioningContext {
	return &cosmosdbProvisioningContext{}
}

func (
	m *module,
) GetEmptyUpdatingParameters() service.UpdatingParameters {
	return &UpdatingParameters{}
}

func (m *module) GetEmptyBindingParameters() service.BindingParameters {
	return &BindingParameters{}
}

func (m *module) GetEmptyBindingContext() service.BindingContext {
	return &cosmosdbBindingContext{}
}

func (m *module) GetEmptyCredentials() service.Credentials {
	return &Credentials{}
}

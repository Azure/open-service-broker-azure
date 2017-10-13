package cosmosdb

import "github.com/Azure/azure-service-broker/pkg/service"

type databaseKind string

const (
	databaseKindMongoDB          databaseKind = "MongoDB"
	databaseKindGlobalDocumentDB databaseKind = "GlobalDocumentDB"
)

// ProvisioningParameters encapsulates CosmosDB-specific provisioning options
type ProvisioningParameters struct {
	Location      string            `json:"location"`
	ResourceGroup string            `json:"resourceGroup"`
	Tags          map[string]string `json:"tags"`
}

type cosmosdbProvisioningContext struct {
	ResourceGroupName        string       `json:"resourceGroup"`
	ARMDeploymentName        string       `json:"armDeployment"`
	DatabaseAccountName      string       `json:"name"`
	DatabaseKind             databaseKind `json:"kind"`
	FullyQualifiedDomainName string       `json:"fullyQualifiedDomainName"`
	ConnectionString         string       `json:"connectionString"`
	PrimaryKey               string       `json:"primaryKey"`
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
	ConnectionString        string `json:"connectionstring,omitempty"`
	URI                     string `json:"uri,omitempty"`
	PrimaryConnectionString string `json:"primaryconnectionstring,omitempty"`
	PrimaryKey              string `json:"primarykey,omitempty"`
}

func (m *cosmosdbProvisioningContext) GetResourceGroupName() string {
	return m.ResourceGroupName
}

func (
	m *module,
) GetEmptyProvisioningParameters() service.ProvisioningParameters {
	return &ProvisioningParameters{}
}

func (m *module) GetEmptyProvisioningContext() service.ProvisioningContext {
	return &cosmosdbProvisioningContext{}
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

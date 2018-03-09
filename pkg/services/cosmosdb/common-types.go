package cosmosdb

type databaseKind string

const (
	databaseKindMongoDB          databaseKind = "MongoDB"
	databaseKindGlobalDocumentDB databaseKind = "GlobalDocumentDB"
)

// ProvisioningParameters encapsulates non-sensitive CosmosDB-specific
// provisioning options
type ProvisioningParameters struct{}

// SecureProvisioningParameters encapsulates sensitive CosmosDB-specific
// provisioning options
type SecureProvisioningParameters struct{}

type cosmosdbInstanceDetails struct {
	ARMDeploymentName        string       `json:"armDeployment"`
	DatabaseAccountName      string       `json:"name"`
	DatabaseKind             databaseKind `json:"kind"`
	FullyQualifiedDomainName string       `json:"fullyQualifiedDomainName"`
}

type cosmosdbSecureInstanceDetails struct {
	ConnectionString string `json:"connectionString"`
	PrimaryKey       string `json:"primaryKey"`
}

// BindingParameters encapsulates non-sensitive CosmosDB-specific binding
// options
type BindingParameters struct {
}

// SecureBindingParameters encapsulates sensitive CosmosDB-specific binding
// options
type SecureBindingParameters struct {
}

type cosmosdbBindingDetails struct {
}

type cosmosdbSecureBindingDetails struct {
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

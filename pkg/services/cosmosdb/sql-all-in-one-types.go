package cosmosdb

type sqlAllInOneInstanceDetails struct {
	cosmosdbInstanceDetails `json:",squash"`
	DatabaseName            string `json:"databaseName"`
}

// cosmosCredentials encapsulates CosmosDB-specific details for connecting via
// a variety of APIs. This excludes MongoDB.
type sqlAPICredentials struct {
	URI                     string `json:"uri,omitempty"`
	PrimaryConnectionString string `json:"primaryConnectionString,omitempty"`
	PrimaryKey              string `json:"primaryKey,omitempty"`
	DatabaseName            string `json:"databaseName"`
	DatabaseID              string `json:"documentdb_database_id"`
	Host                    string `json:"documentdb_host_endpoint"`
	MasterKey               string `json:"documentdb_master_key"`
}

type databaseCreationRequest struct {
	ID string `json:"id"`
}

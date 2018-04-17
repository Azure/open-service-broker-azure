package cosmosdb

// mongoCredentials encapsulates CosmosDB-specific connection details and
// credentials for connecting with the MongoDB API.
type mongoCredentials struct {
	Host             string `json:"host,omitempty"`
	Port             int    `json:"port,omitempty"`
	Username         string `json:"username,omitempty"`
	Password         string `json:"password,omitempty"`
	ConnectionString string `json:"connectionString,omitempty"`
	URI              string `json:"uri,omitempty"`
}

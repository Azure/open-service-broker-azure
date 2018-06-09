// +build experimental

package cosmosdb

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

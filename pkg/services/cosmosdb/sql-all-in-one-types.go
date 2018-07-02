package cosmosdb

import "github.com/Azure/open-service-broker-azure/pkg/service"

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

// GetEmptyInstanceDetails returns an "empty" service-specific object that
// can be populated with data during unmarshaling of JSON to an Instance
func (
	s *sqlAllInOneManager,
) GetEmptyInstanceDetails() service.InstanceDetails {
	return &sqlAllInOneInstanceDetails{}
}

// GetEmptyBindingDetails returns an "empty" service-specific object that
// can be populated with data during unmarshaling of JSON to a Binding
func (s *sqlAllInOneManager) GetEmptyBindingDetails() service.BindingDetails {
	return nil
}

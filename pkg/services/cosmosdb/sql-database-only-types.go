package cosmosdb

import "github.com/Azure/open-service-broker-azure/pkg/service"

type sqlDatabaseOnlyInstanceDetails struct {
	DatabaseName string `json:"databaseName"`
}

// GetEmptyInstanceDetails returns an "empty" service-specific object that
// can be populated with data during unmarshaling of JSON to an Instance
func (
	s *sqlDatabaseManager,
) GetEmptyInstanceDetails() service.InstanceDetails {
	return &sqlDatabaseOnlyInstanceDetails{}
}

// GetEmptyBindingDetails returns an "empty" service-specific object that
// can be populated with data during unmarshaling of JSON to a Binding
func (s *sqlDatabaseManager) GetEmptyBindingDetails() service.BindingDetails {
	return nil
}

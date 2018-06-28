package cosmosdb

import "github.com/Azure/open-service-broker-azure/pkg/service"

type cosmosdbInstanceDetails struct {
	ARMDeploymentName        string `json:"armDeployment"`
	DatabaseAccountName      string `json:"name"`
	FullyQualifiedDomainName string `json:"fullyQualifiedDomainName"`
	IPFilters                string `json:"ipFilters"`
	ConnectionString         string `json:"connectionString"`
	PrimaryKey               string `json:"primaryKey"`
}

// cosmosCredentials encapsulates CosmosDB-specific details for connecting via
// a variety of APIs. This excludes MongoDB.
type cosmosCredentials struct {
	URI                     string `json:"uri"`
	PrimaryConnectionString string `json:"primaryConnectionString"`
	PrimaryKey              string `json:"primaryKey"`
}

// GetEmptyInstanceDetails returns an "empty" service-specific object that
// can be populated with data during unmarshaling of JSON to an Instance
func (
	c *cosmosAccountManager,
) GetEmptyInstanceDetails() service.InstanceDetails {
	return &cosmosdbInstanceDetails{}
}

// GetEmptyBindingDetails returns an "empty" service-specific object that
// can be populated with data during unmarshaling of JSON to a Binding
func (c *cosmosAccountManager) GetEmptyBindingDetails() service.BindingDetails {
	return nil
}

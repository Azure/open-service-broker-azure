package cosmosdb

import (
	"github.com/Azure/open-service-broker-azure/pkg/service"
)

func (
	s *sqlAllInOneManager,
) getProvisionParametersSchema() map[string]service.ParameterSchema {
	p := s.cosmosAccountManager.getProvisionParametersSchema()

	p["storageCapacity"] = &service.SimpleParameterSchema{
		Type: "string",
		Description: "Represents the maximum storage size of the collection." +
			" Can be either `fixed` (10 GB) or `unlimited`. If unspecified " +
			"defaults to `unlimited`",
		AllowedValues: []string{"", "fixed", "unlimited"},
		Default:       "unlimited",
	}

	p["partitionKey"] = &service.SimpleParameterSchema{
		Type: "string",
		Description: "The Partition Key is used to automatically partition" +
			"data among multiple servers for scalability. Choose a JSON" +
			" property name that has a wide range of values and is likely" +
			"  to have evenly distributed access patterns. Required if" +
			" storage is unlimited. Otherwise ignored.",
	}

	p["throughput"] = &service.SimpleParameterSchema{
		Type: "integer",
		Description: "Each collection can be provisioned throughput in " +
			"Request Units per second (RU/s). 1 RU corresponds to the " +
			"throughput of a read of a 1 KB document. Maximum value is ",
	}

	p["uniqueKeys"] = &service.ArrayParameterSchema{
		ItemsSchema: &service.SimpleParameterSchema{
			Type: "string",
		},
		Description: "Unique keys provide developers with the ability to " +
			"add a layer of data integrity to their database. By creating a " +
			"unique key policy when a container is created, you ensure the " +
			"uniqueness of one or more values per partition key.",
	}
	return p
}

// type cosmosdbInstanceDetails struct {
// 	ARMDeploymentName        string `json:"armDeployment"`
// 	DatabaseAccountName      string `json:"name"`
// 	FullyQualifiedDomainName string `json:"fullyQualifiedDomainName"`
// 	IPFilters                string `json:"ipFilters"`
// }

type sqlAllInOneInstanceDetails struct {
	ARMDeploymentName        string `json:"armDeployment"`
	DatabaseAccountName      string `json:"name"`
	FullyQualifiedDomainName string `json:"fullyQualifiedDomainName"`
	IPFilters                string `json:"ipFilters"`
	DatabaseName             string `json:"databaseName"`
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

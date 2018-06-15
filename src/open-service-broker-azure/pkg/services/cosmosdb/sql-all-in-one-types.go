// +build experimental

package cosmosdb

import (
	"open-service-broker-azure/pkg/service"
)

func (
	s *sqlAllInOneManager,
) getProvisionParametersSchema() service.InputParametersSchema {
	return s.cosmosAccountManager.getProvisionParametersSchema()
}

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

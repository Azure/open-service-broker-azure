package cosmosdb

import (
	"fmt"

	"github.com/Azure/open-service-broker-azure/pkg/service"
)

func (s *serviceManager) ValidateBindingParameters(
	bindingParameters service.BindingParameters,
) error {
	// There are no parameters for binding to CosmosDB, so there is nothing
	// to validate
	return nil
}

func (s *serviceManager) Bind(
	instance service.Instance,
	_ service.BindingParameters,
) (service.BindingDetails, service.Credentials, error) {
	dt, ok := instance.Details.(*cosmosdbInstanceDetails)
	if !ok {
		return nil, nil, fmt.Errorf(
			"error casting instance.Details as *cosmosdbInstanceDetails",
		)
	}
	if dt.DatabaseKind == databaseKindMongoDB {
		cosmosDBCredentials := &Credentials{
			Host: dt.FullyQualifiedDomainName,
			Port: 10255,
			// Username is the same as the database account name
			Username:         dt.DatabaseAccountName,
			Password:         dt.PrimaryKey,
			ConnectionString: dt.ConnectionString,
		}
		return &cosmosdbBindingDetails{},
			cosmosDBCredentials,
			nil
	}
	cosmosDBCredentials := &Credentials{
		URI:                     dt.FullyQualifiedDomainName,
		PrimaryKey:              dt.PrimaryKey,
		PrimaryConnectionString: dt.ConnectionString,
	}
	return &cosmosdbBindingDetails{},
		cosmosDBCredentials,
		nil
}

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
	service.Instance,
	service.BindingParameters,
) (service.BindingDetails, error) {
	return &cosmosdbBindingDetails{}, nil
}

func (s *serviceManager) GetCredentials(
	instance service.Instance,
	_ service.Binding,
) (service.Credentials, error) {
	dt, ok := instance.Details.(*cosmosdbInstanceDetails)
	if !ok {
		return nil, fmt.Errorf(
			"error casting instance.Details as *cosmosdbInstanceDetails",
		)
	}
	if dt.DatabaseKind == databaseKindMongoDB {
		return &Credentials{
			Host: dt.FullyQualifiedDomainName,
			Port: 10255,
			// Username is the same as the database account name
			Username:         dt.DatabaseAccountName,
			Password:         dt.PrimaryKey,
			ConnectionString: dt.ConnectionString,
		}, nil
	}
	return &Credentials{
		URI:                     dt.FullyQualifiedDomainName,
		PrimaryKey:              dt.PrimaryKey,
		PrimaryConnectionString: dt.ConnectionString,
	}, nil
}

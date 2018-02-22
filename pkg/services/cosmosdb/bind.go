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
	service.SecureBindingParameters,
) (service.BindingDetails, service.SecureBindingDetails, error) {
	return &cosmosdbBindingDetails{}, &cosmosdbSecureBindingDetails{}, nil
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
	sdt, ok := instance.SecureDetails.(*cosmosdbSecureInstanceDetails)
	if !ok {
		return nil, fmt.Errorf(
			"error casting instance.SecureDetails as *cosmosdbSecureInstanceDetails",
		)
	}
	if dt.DatabaseKind == databaseKindMongoDB {
		return &Credentials{
			Host: dt.FullyQualifiedDomainName,
			Port: 10255,
			// Username is the same as the database account name
			Username:         dt.DatabaseAccountName,
			Password:         sdt.PrimaryKey,
			ConnectionString: sdt.ConnectionString,
		}, nil
	}
	return &Credentials{
		URI:                     dt.FullyQualifiedDomainName,
		PrimaryKey:              sdt.PrimaryKey,
		PrimaryConnectionString: sdt.ConnectionString,
	}, nil
}

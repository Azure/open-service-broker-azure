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
) (service.BindingContext, service.Credentials, error) {
	pc, ok := instance.ProvisioningContext.(*cosmosdbProvisioningContext)
	if !ok {
		return nil, nil, fmt.Errorf(
			"error casting instance.ProvisioningContext as " +
				"*cosmosdbProvisioningContext",
		)
	}
	if pc.DatabaseKind == databaseKindMongoDB {
		cosmosDBCredentials := &Credentials{
			Host: pc.FullyQualifiedDomainName,
			Port: 10255,
			// Username is the same as the database account name
			Username:         pc.DatabaseAccountName,
			Password:         pc.PrimaryKey,
			ConnectionString: pc.ConnectionString,
		}
		return &cosmosdbBindingContext{},
			cosmosDBCredentials,
			nil
	}
	cosmosDBCredentials := &Credentials{
		URI:                     pc.FullyQualifiedDomainName,
		PrimaryKey:              pc.PrimaryKey,
		PrimaryConnectionString: pc.ConnectionString,
	}
	return &cosmosdbBindingContext{},
		cosmosDBCredentials,
		nil

}

package cosmosdb

import (
	"fmt"

	"github.com/Azure/azure-service-broker/pkg/service"
)

func (m *module) ValidateBindingParameters(
	bindingParameters service.BindingParameters,
) error {
	// There are no parameters for binding to CosmosDB, so there is nothing
	// to validate
	return nil
}

func (m *module) Bind(
	provisioningContext service.ProvisioningContext,
	bindingParameters service.BindingParameters,
) (service.BindingContext, service.Credentials, error) {
	pc, ok := provisioningContext.(*cosmosdbProvisioningContext)
	if !ok {
		return nil, nil, fmt.Errorf(
			"error casting provisioningContext as *cosmosdbProvisioningContext",
		)
	}
	if pc.DatabaseKind == databaseKindMongoDB {
		cosmosDBCredentials := &cosmosdbCredentials{
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
	cosmosDBCredentials := &cosmosdbCredentials{
		URI:                     pc.FullyQualifiedDomainName,
		PrimaryKey:              pc.PrimaryKey,
		PrimaryConnectionString: pc.ConnectionString,
	}
	return &cosmosdbBindingContext{},
		cosmosDBCredentials,
		nil

}

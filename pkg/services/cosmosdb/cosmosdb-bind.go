package cosmosdb

import (
	"fmt"

	"github.com/Azure/open-service-broker-azure/pkg/service"
)

func (c *cosmosManager) ValidateBindingParameters(
	service.BindingParameters,
	service.SecureBindingParameters,
) error {
	// There are no parameters for binding to CosmosDB, so there is nothing
	// to validate
	return nil
}

func (c *cosmosManager) Bind(
	service.Instance,
	service.BindingParameters,
	service.SecureBindingParameters,
) (service.BindingDetails, service.SecureBindingDetails, error) {
	return &cosmosdbBindingDetails{}, &cosmosdbSecureBindingDetails{}, nil
}

func (c *cosmosManager) GetCredentials(
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
	return &CosmosCredentials{
		URI:                     dt.FullyQualifiedDomainName,
		PrimaryKey:              sdt.PrimaryKey,
		PrimaryConnectionString: sdt.ConnectionString,
	}, nil
}

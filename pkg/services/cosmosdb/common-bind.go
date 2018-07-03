package cosmosdb

import (
	"github.com/Azure/open-service-broker-azure/pkg/service"
)

func (c *cosmosAccountManager) Bind(
	service.Instance,
	service.BindingParameters,
) (service.BindingDetails, error) {
	return nil, nil
}

func (c *cosmosAccountManager) GetCredentials(
	instance service.Instance,
	_ service.Binding,
) (service.Credentials, error) {
	dt := instance.Details.(*cosmosdbInstanceDetails)
	return cosmosCredentials{
		URI:                     dt.FullyQualifiedDomainName,
		PrimaryKey:              string(dt.PrimaryKey),
		PrimaryConnectionString: string(dt.ConnectionString),
	}, nil
}

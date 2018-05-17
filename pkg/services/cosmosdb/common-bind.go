package cosmosdb

import (
	"github.com/Azure/open-service-broker-azure/pkg/service"
)

func (c *cosmosAccountManager) Bind(
	service.Instance,
	service.BindingParameters,
	service.SecureBindingParameters,
) (service.BindingDetails, service.SecureBindingDetails, error) {
	return nil, nil, nil
}

func (c *cosmosAccountManager) GetCredentials(
	instance service.Instance,
	_ service.Binding,
) (service.Credentials, error) {
	dt := cosmosdbInstanceDetails{}
	if err := service.GetStructFromMap(instance.Details, &dt); err != nil {
		return nil, err
	}
	sdt := cosmosdbSecureInstanceDetails{}
	if err := service.GetStructFromMap(instance.SecureDetails, &sdt); err != nil {
		return nil, err
	}
	return cosmosCredentials{
		URI:                     dt.FullyQualifiedDomainName,
		PrimaryKey:              sdt.PrimaryKey,
		PrimaryConnectionString: sdt.ConnectionString,
	}, nil
}

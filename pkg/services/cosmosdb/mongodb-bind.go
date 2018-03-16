package cosmosdb

import (
	"fmt"

	"github.com/Azure/open-service-broker-azure/pkg/service"
)

func (m *mongoAccountManager) ValidateBindingParameters(
	service.BindingParameters,
	service.SecureBindingParameters,
) error {
	// There are no parameters for binding to CosmosDB, so there is nothing
	// to validate
	return nil
}

func (m *mongoAccountManager) Bind(
	service.Instance,
	service.BindingParameters,
	service.SecureBindingParameters,
) (service.BindingDetails, service.SecureBindingDetails, error) {
	return nil, nil, nil
}

func (m *mongoAccountManager) GetCredentials(
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
	return &MongoCredentials{
		Host: dt.FullyQualifiedDomainName,
		Port: 10255,
		// Username is the same as the database account name
		Username:         dt.DatabaseAccountName,
		Password:         sdt.PrimaryKey,
		ConnectionString: sdt.ConnectionString,
	}, nil
}

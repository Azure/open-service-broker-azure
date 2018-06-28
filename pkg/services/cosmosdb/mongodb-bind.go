package cosmosdb

import (
	"github.com/Azure/open-service-broker-azure/pkg/service"
)

func (m *mongoAccountManager) GetCredentials(
	instance service.Instance,
	_ service.Binding,
) (service.Credentials, error) {
	dt := instance.Details.(*cosmosdbInstanceDetails)
	return mongoCredentials{
		Host: dt.FullyQualifiedDomainName,
		Port: 10255,
		// Username is the same as the database account name
		Username:         dt.DatabaseAccountName,
		Password:         dt.PrimaryKey,
		ConnectionString: dt.ConnectionString,
		URI:              dt.ConnectionString,
	}, nil
}

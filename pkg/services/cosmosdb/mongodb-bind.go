// +build experimental

package cosmosdb

import (
	"github.com/Azure/open-service-broker-azure/pkg/service"
)

func (m *mongoAccountManager) GetCredentials(
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
	return mongoCredentials{
		Host: dt.FullyQualifiedDomainName,
		Port: 10255,
		// Username is the same as the database account name
		Username:         dt.DatabaseAccountName,
		Password:         sdt.PrimaryKey,
		ConnectionString: sdt.ConnectionString,
		URI:              sdt.ConnectionString,
	}, nil
}

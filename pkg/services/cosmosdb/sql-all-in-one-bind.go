// +build experimental

package cosmosdb

import (
	"github.com/Azure/open-service-broker-azure/pkg/service"
)

func (s *sqlAllInOneManager) GetCredentials(
	instance service.Instance,
	_ service.Binding,
) (service.Credentials, error) {
	dt := sqlAllInOneInstanceDetails{}
	if err := service.GetStructFromMap(instance.Details, &dt); err != nil {
		return nil, err
	}
	sdt := cosmosdbSecureInstanceDetails{}
	if err := service.GetStructFromMap(instance.SecureDetails, &sdt); err != nil {
		return nil, err
	}
	return sqlAPICredentials{
		URI:                     dt.FullyQualifiedDomainName,
		PrimaryKey:              sdt.PrimaryKey,
		PrimaryConnectionString: sdt.ConnectionString,
		DatabaseName:            dt.DatabaseName,
		DatabaseID:              dt.DatabaseName,
		Host:                    dt.FullyQualifiedDomainName,
		MasterKey:               sdt.PrimaryKey,
	}, nil
}

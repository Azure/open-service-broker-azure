package cosmosdb

import (
	"github.com/Azure/open-service-broker-azure/pkg/service"
)

func (s *sqlAllInOneManager) GetCredentials(
	instance service.Instance,
	_ service.Binding,
) (service.Credentials, error) {
	dt := instance.Details.(*sqlAllInOneInstanceDetails)
	return sqlAPICredentials{
		URI:                     dt.FullyQualifiedDomainName,
		PrimaryKey:              string(dt.PrimaryKey),
		PrimaryConnectionString: string(dt.ConnectionString),
		DatabaseName:            dt.DatabaseName,
		DatabaseID:              dt.DatabaseName,
		Host:                    dt.FullyQualifiedDomainName,
		MasterKey:               string(dt.PrimaryKey),
	}, nil
}

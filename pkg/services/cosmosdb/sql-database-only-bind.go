package cosmosdb

import (
	"github.com/Azure/open-service-broker-azure/pkg/service"
)

func (s *sqlDatabaseManager) GetCredentials(
	instance service.Instance,
	_ service.Binding,
) (service.Credentials, error) {
	dt := instance.Details.(*sqlDatabaseOnlyInstanceDetails)
	pdt := instance.Parent.Details.(*cosmosdbInstanceDetails)
	return sqlAPICredentials{
		URI:                     pdt.FullyQualifiedDomainName,
		PrimaryKey:              pdt.PrimaryKey,
		PrimaryConnectionString: pdt.ConnectionString,
		DatabaseName:            dt.DatabaseName,
		DatabaseID:              dt.DatabaseName,
		Host:                    pdt.FullyQualifiedDomainName,
		MasterKey:               pdt.PrimaryKey,
	}, nil
}

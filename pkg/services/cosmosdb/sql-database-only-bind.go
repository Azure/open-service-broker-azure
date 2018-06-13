// +build experimental

package cosmosdb

import (
	"github.com/Azure/open-service-broker-azure/pkg/service"
)

func (s *sqlDatabaseManager) GetCredentials(
	instance service.Instance,
	_ service.Binding,
) (service.Credentials, error) {
	dt := sqlDatabaseOnlyInstanceDetails{}
	if err := service.GetStructFromMap(instance.Details, &dt); err != nil {
		return nil, err
	}
	pdt := cosmosdbInstanceDetails{}
	if err := service.GetStructFromMap(instance.Details, &pdt); err != nil {
		return nil, err
	}
	psdt := cosmosdbSecureInstanceDetails{}
	if err := service.GetStructFromMap(
		instance.Parent.SecureDetails,
		&psdt,
	); err != nil {
		return nil, err
	}
	return sqlAPICredentials{
		URI:                     pdt.FullyQualifiedDomainName,
		PrimaryKey:              psdt.PrimaryKey,
		PrimaryConnectionString: psdt.ConnectionString,
		DatabaseName:            dt.DatabaseName,
		DatabaseID:              dt.DatabaseName,
		Host:                    pdt.FullyQualifiedDomainName,
		MasterKey:               psdt.PrimaryKey,
	}, nil
}

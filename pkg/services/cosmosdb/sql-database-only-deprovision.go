package cosmosdb

import (
	"context"

	"github.com/Azure/open-service-broker-azure/pkg/service"
)

func (s *sqlDatabaseManager) GetDeprovisioner(
	service.Plan,
) (service.Deprovisioner, error) {
	return service.NewDeprovisioner(
		service.NewDeprovisioningStep("deleteDatabase", s.deleteDatabase),
	)
}

func (s *sqlDatabaseManager) deleteDatabase(
	_ context.Context,
	instance service.Instance,
) (service.InstanceDetails, error) {
	dt := instance.Details.(*sqlDatabaseOnlyInstanceDetails)
	pdt := instance.Parent.Details.(*cosmosdbInstanceDetails)
	err := deleteDatabase(pdt.DatabaseAccountName, dt.DatabaseName, pdt.PrimaryKey)
	if err != nil {
		return nil, err
	}
	return instance.Details, nil
}

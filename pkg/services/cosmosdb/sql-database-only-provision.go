package cosmosdb

import (
	"context"

	"github.com/Azure/open-service-broker-azure/pkg/service"
	uuid "github.com/satori/go.uuid"
)

func (s *sqlDatabaseManager) GetProvisioner(
	service.Plan,
) (service.Provisioner, error) {
	return service.NewProvisioner(
		service.NewProvisioningStep("preProvision", s.preProvision),
		service.NewProvisioningStep("createDatabase", s.createDatabase),
	)
}

func (s *sqlDatabaseManager) preProvision(
	_ context.Context,
	_ service.Instance,
) (service.InstanceDetails, error) {
	return &sqlDatabaseOnlyInstanceDetails{
		DatabaseName: uuid.NewV4().String(),
	}, nil
}

func (s *sqlDatabaseManager) createDatabase(
	_ context.Context,
	instance service.Instance,
) (service.InstanceDetails, error) {
	dt := instance.Details.(*sqlDatabaseOnlyInstanceDetails)
	pdt := instance.Parent.Details.(*cosmosdbInstanceDetails)
	err := createDatabase(pdt.DatabaseAccountName, dt.DatabaseName, pdt.PrimaryKey)
	if err != nil {
		return nil, err
	}
	return instance.Details, nil
}

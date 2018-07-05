package cosmosdb

import (
	"context"

	"github.com/Azure/open-service-broker-azure/pkg/service"
)

func (s *sqlAllInOneManager) GetDeprovisioner(
	service.Plan,
) (service.Deprovisioner, error) {
	return service.NewDeprovisioner(
		service.NewDeprovisioningStep("deleteARMDeployment", s.deleteARMDeployment),
		service.NewDeprovisioningStep(
			"deleteDatabase", s.deleteDatabase,
		),
		service.NewDeprovisioningStep(
			"deleteCosmosDBServer", s.deleteCosmosDBAccount,
		),
	)
}

func (s *sqlAllInOneManager) deleteARMDeployment(
	_ context.Context,
	instance service.Instance,
) (service.InstanceDetails, error) {
	dt := instance.Details.(*sqlAllInOneInstanceDetails)
	if err := deleteARMDeployment(
		s.armDeployer,
		instance.ProvisioningParameters,
		&dt.cosmosdbInstanceDetails,
	); err != nil {
		return nil, err
	}
	return dt, nil
}

func (s *sqlAllInOneManager) deleteDatabase(
	_ context.Context,
	instance service.Instance,
) (service.InstanceDetails, error) {
	dt := instance.Details.(*sqlAllInOneInstanceDetails)
	err := deleteDatabase(
		dt.DatabaseAccountName,
		dt.DatabaseName,
		string(dt.PrimaryKey),
	)
	if err != nil {
		return nil, err
	}
	return instance.Details, nil
}

func (s *sqlAllInOneManager) deleteCosmosDBAccount(
	ctx context.Context,
	instance service.Instance,
) (service.InstanceDetails, error) {
	dt := instance.Details.(*sqlAllInOneInstanceDetails)
	if err := deleteCosmosDBAccount(
		ctx,
		s.databaseAccountsClient,
		instance.ProvisioningParameters,
		&dt.cosmosdbInstanceDetails,
	); err != nil {
		return nil, err
	}
	return instance.Details, nil
}

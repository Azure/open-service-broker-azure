package cosmosdb

import (
	"context"
	"fmt"

	"github.com/Azure/open-service-broker-azure/pkg/service"
)

func (
	s *sqlAllInOneManager,
) GetUpdater(service.Plan) (service.Updater, error) {
	return service.NewUpdater(
		service.NewUpdatingStep("updateARMTemplate", s.updateARMTemplate),
		service.NewUpdatingStep("waitForReadLocationsReadyInUpdate", s.waitForReadLocationsReadyInUpdate), //nolint: lll
	)
}

func (s *sqlAllInOneManager) updateARMTemplate(
	_ context.Context,
	instance service.Instance,
) (service.InstanceDetails, error) {
	dt := instance.Details.(*sqlAllInOneInstanceDetails)
	if err := s.cosmosAccountManager.updateDeployment(
		instance.ProvisioningParameters,
		instance.UpdatingParameters,
		&dt.cosmosdbInstanceDetails,
		"GlobalDocumentDB",
		"",
		map[string]string{
			"defaultExperience": "DocumentDB",
		},
	); err != nil {
		return nil, fmt.Errorf("error deploying ARM template: %s", err)
	}
	return instance.Details, nil
}

// This function is the same as `s.waitForReadLocationsReady` except that
// it uses `readRegions` array in updating parameter.
func (s *sqlAllInOneManager) waitForReadLocationsReadyInUpdate(
	ctx context.Context,
	instance service.Instance,
) (service.InstanceDetails, error) {
	dt := instance.Details.(*sqlAllInOneInstanceDetails)
	resourceGroupName := instance.ProvisioningParameters.GetString("resourceGroup")
	accountName := dt.DatabaseAccountName
	databaseAccountClient := s.databaseAccountsClient

	err := pollingUntilReadLocationsReady(
		ctx,
		resourceGroupName,
		accountName,
		databaseAccountClient,
		instance.ProvisioningParameters.GetString("location"),
		instance.UpdatingParameters.GetStringArray("readRegions"),
		false,
	)
	if err != nil {
		return nil, err
	}
	return dt, nil
}

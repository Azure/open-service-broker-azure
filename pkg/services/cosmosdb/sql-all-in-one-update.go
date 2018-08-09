package cosmosdb

import (
	"context"
	"fmt"

	"github.com/Azure/open-service-broker-azure/pkg/service"
)

func (
	s *sqlAllInOneManager,
) ValidateUpdatingParameters(instance service.Instance) error {
	return validateReadLocations(
		"sql all in one update",
		instance.UpdatingParameters.GetStringArray("readRegions"),
	)
}

func (
	s *sqlAllInOneManager,
) GetUpdater(service.Plan) (service.Updater, error) {
	return service.NewUpdater(
		service.NewUpdatingStep("updateReadLocations", s.updateReadLocations),
		service.NewUpdatingStep("waitForReadLocationsReady", s.waitForReadLocationsReady), //nolint: lll
		service.NewUpdatingStep("updateARMTemplate", s.updateARMTemplate),
	)
}

func (s *sqlAllInOneManager) updateReadLocations(
	_ context.Context,
	instance service.Instance,
) (service.InstanceDetails, error) {
	dt := instance.Details.(*sqlAllInOneInstanceDetails)
	if err := s.cosmosAccountManager.updateReadLocations(
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

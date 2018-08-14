package cosmosdb

import (
	"context"
	"fmt"

	"github.com/Azure/open-service-broker-azure/pkg/service"
)

func (
	s *sqlAccountManager,
) ValidateUpdatingParameters(instance service.Instance) error {
	return validateReadLocations(
		"sql account update",
		instance.UpdatingParameters.GetStringArray("readRegions"),
	)
}

func (
	s *sqlAccountManager,
) GetUpdater(service.Plan) (service.Updater, error) {
	return service.NewUpdater(
		service.NewUpdatingStep("updateReadLocations", s.updateReadLocations),
		service.NewUpdatingStep("waitForReadLocationsReadyInUpdate", s.waitForReadLocationsReadyInUpdate), //nolint: lll
		service.NewUpdatingStep("updateARMTemplate", s.updateARMTemplate),
	)
}

func (s *sqlAccountManager) updateReadLocations(
	_ context.Context,
	instance service.Instance,
) (service.InstanceDetails, error) {
	if err := s.cosmosAccountManager.updateReadLocations(
		instance.ProvisioningParameters,
		instance.UpdatingParameters,
		instance.Details.(*cosmosdbInstanceDetails),
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

func (s *sqlAccountManager) updateARMTemplate(
	_ context.Context,
	instance service.Instance,
) (service.InstanceDetails, error) {
	if err := s.cosmosAccountManager.updateDeployment(
		instance.ProvisioningParameters,
		instance.UpdatingParameters,
		instance.Details.(*cosmosdbInstanceDetails),
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

package cosmosdb

import (
	"context"
	"fmt"

	"github.com/Azure/open-service-broker-azure/pkg/service"
)

func (
	s *sqlAccountManager,
) GetUpdater(service.Plan) (service.Updater, error) {
	return service.NewUpdater(
		service.NewUpdatingStep("updateARMTemplate", s.updateARMTemplate),
		service.NewUpdatingStep("waitForReadLocationsReadyInUpdate", s.waitForReadLocationsReadyInUpdate), //nolint: lll
	)
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

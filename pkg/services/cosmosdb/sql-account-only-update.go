package cosmosdb

import (
	"context"
	"fmt"

	"github.com/Azure/open-service-broker-azure/pkg/service"
)

func (s *sqlAccountManager) ValidateUpdatingParameters(
	instance service.Instance,
) error {
	return nil
}

func (s *sqlAccountManager) GetUpdater(service.Plan) (service.Updater, error) {
	// There isn't a need to do any "pre-provision here. just the update step"
	return service.NewUpdater(
		service.NewUpdatingStep("updateARMTemplate", s.updateARMTemplate),
	)
}

func (s *sqlAccountManager) updateARMTemplate(
	_ context.Context,
	instance service.Instance,
) (service.InstanceDetails, error) {
	err := s.cosmosAccountManager.updateDeployment(
		instance.ProvisioningParameters,
		instance.UpdatingParameters,
		instance.Details.(*cosmosdbInstanceDetails),
		"GlobalDocumentDB",
		"",
		map[string]string{
			"defaultExperience": "DocumentDB",
		},
	)
	if err != nil {
		return nil, fmt.Errorf("error deploying ARM template: %s", err)
	}
	return instance.Details, nil
}

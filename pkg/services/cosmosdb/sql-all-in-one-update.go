package cosmosdb

import (
	"context"
	"fmt"

	"github.com/Azure/open-service-broker-azure/pkg/service"
)

func (s *sqlAllInOneManager) ValidateUpdatingParameters(
	instance service.Instance,
) error {
	return nil
}

func (s *sqlAllInOneManager) GetUpdater(service.Plan) (service.Updater, error) {
	// There isn't a need to do any "pre-provision here. just the update step"
	return service.NewUpdater(
		service.NewUpdatingStep("updateARMTemplate", s.updateARMTemplate),
	)
}

func (s *sqlAllInOneManager) updateARMTemplate(
	_ context.Context,
	instance service.Instance,
) (service.InstanceDetails, error) {
	dt := instance.Details.(*sqlAllInOneInstanceDetails)
	err := s.cosmosAccountManager.updateDeployment(
		instance.ProvisioningParameters,
		instance.UpdatingParameters,
		&dt.cosmosdbInstanceDetails,
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

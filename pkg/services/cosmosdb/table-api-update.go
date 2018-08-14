package cosmosdb

import (
	"context"
	"fmt"

	"github.com/Azure/open-service-broker-azure/pkg/service"
)

func (
	t *tableAccountManager,
) ValidateUpdatingParameters(instance service.Instance) error {
	return validateReadLocations(
		"table account update",
		instance.UpdatingParameters.GetStringArray("readRegions"),
	)
}

func (
	t *tableAccountManager,
) GetUpdater(service.Plan) (service.Updater, error) {
	return service.NewUpdater(
		service.NewUpdatingStep("updateReadLocations", t.updateReadLocations),
		service.NewUpdatingStep("waitForReadLocationsReadyInUpdate", t.waitForReadLocationsReadyInUpdate), //nolint: lll
		service.NewUpdatingStep("updateARMTemplate", t.updateARMTemplate),
	)
}

func (t *tableAccountManager) updateReadLocations(
	_ context.Context,
	instance service.Instance,
) (service.InstanceDetails, error) {
	if err := t.cosmosAccountManager.updateReadLocations(
		instance.ProvisioningParameters,
		instance.UpdatingParameters,
		instance.Details.(*cosmosdbInstanceDetails),
		"GlobalDocumentDB",
		"EnableTable",
		map[string]string{
			"defaultExperience": "Table",
		},
	); err != nil {
		return nil, fmt.Errorf("error deploying ARM template: %s", err)
	}
	return instance.Details, nil
}

func (t *tableAccountManager) updateARMTemplate(
	_ context.Context,
	instance service.Instance,
) (service.InstanceDetails, error) {
	if err := t.cosmosAccountManager.updateDeployment(
		instance.ProvisioningParameters,
		instance.UpdatingParameters,
		instance.Details.(*cosmosdbInstanceDetails),
		"GlobalDocumentDB",
		"EnableTable",
		map[string]string{
			"defaultExperience": "Table",
		},
	); err != nil {
		return nil, fmt.Errorf("error deploying ARM template: %s", err)
	}
	return instance.Details, nil
}

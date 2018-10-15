package cosmosdb

import (
	"context"
	"fmt"

	"github.com/Azure/open-service-broker-azure/pkg/service"
)

func (
	t *tableAccountManager,
) GetUpdater(service.Plan) (service.Updater, error) {
	return service.NewUpdater(
		service.NewUpdatingStep("updateARMTemplate", t.updateARMTemplate),
		service.NewUpdatingStep("waitForReadLocationsReadyInUpdate", t.waitForReadLocationsReadyInUpdate), //nolint: lll
	)
}

func (t *tableAccountManager) updateARMTemplate(
	_ context.Context,
	instance service.Instance,
) (service.InstanceDetails, error) {
	if err := t.cosmosAccountManager.updateDeployment(
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

package cosmosdb

import (
	"context"
	"fmt"

	"github.com/Azure/open-service-broker-azure/pkg/service"
)

func (t *tableAccountManager) ValidateUpdatingParameters(
	instance service.Instance,
) error {
	return nil
}

func (t *tableAccountManager) GetUpdater(
	service.Plan,
) (service.Updater, error) {
	// There isn't a need to do any "pre-provision here. just the update step"
	return service.NewUpdater(
		service.NewUpdatingStep("updateARMTemplate", t.updateARMTemplate),
	)
}

func (t *tableAccountManager) updateARMTemplate(
	_ context.Context,
	instance service.Instance,
) (service.InstanceDetails, error) {
	err := t.cosmosAccountManager.updateDeployment(
		instance.ProvisioningParameters,
		instance.UpdatingParameters,
		instance.Details.(*cosmosdbInstanceDetails),
		"GlobalDocumentDB",
		"EnableTable",
		map[string]string{
			"defaultExperience": "Table",
		},
	)
	if err != nil {
		return nil, fmt.Errorf("error deploying ARM template: %s", err)
	}
	return instance.Details, nil
}

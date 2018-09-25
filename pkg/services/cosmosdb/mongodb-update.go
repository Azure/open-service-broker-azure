package cosmosdb

import (
	"context"
	"fmt"

	"github.com/Azure/open-service-broker-azure/pkg/service"
)

func (
	m *mongoAccountManager,
) GetUpdater(service.Plan) (service.Updater, error) {
	return service.NewUpdater(
		service.NewUpdatingStep("updateARMTemplate", m.updateARMTemplate),
		service.NewUpdatingStep("waitForReadLocationsReadyInUpdate", m.waitForReadLocationsReadyInUpdate), //nolint: lll
	)
}

func (m *mongoAccountManager) updateARMTemplate(
	_ context.Context,
	instance service.Instance,
) (service.InstanceDetails, error) {
	if err := m.cosmosAccountManager.updateDeployment(
		instance.ProvisioningParameters,
		instance.UpdatingParameters,
		instance.Details.(*cosmosdbInstanceDetails),
		"MongoDB",
		"",
		map[string]string{},
	); err != nil {
		return nil, fmt.Errorf("error deploying ARM template: %s", err)
	}
	return instance.Details, nil
}

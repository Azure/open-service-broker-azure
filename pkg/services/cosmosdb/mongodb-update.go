package cosmosdb

import (
	"context"
	"fmt"

	"github.com/Azure/open-service-broker-azure/pkg/service"
)

func (
	m *mongoAccountManager,
) ValidateUpdatingParameters(instance service.Instance) error {
	return validateReadLocations(
		"mongo account update",
		instance.UpdatingParameters.GetStringArray("readRegions"),
	)
}

func (
	m *mongoAccountManager,
) GetUpdater(service.Plan) (service.Updater, error) {
	return service.NewUpdater(
		service.NewUpdatingStep("updateReadLocations", m.updateReadLocations),
		service.NewUpdatingStep("waitForReadLocationsReadyInUpdate", m.waitForReadLocationsReadyInUpdate), //nolint: lll
		service.NewUpdatingStep("updateARMTemplate", m.updateARMTemplate),
	)
}

func (m *mongoAccountManager) updateReadLocations(
	_ context.Context,
	instance service.Instance,
) (service.InstanceDetails, error) {
	if err := m.cosmosAccountManager.updateReadLocations(
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

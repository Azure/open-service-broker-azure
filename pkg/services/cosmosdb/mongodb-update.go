package cosmosdb

import (
	"context"
	"fmt"

	"github.com/Azure/open-service-broker-azure/pkg/service"
)

func (m *mongoAccountManager) ValidateUpdatingParameters(
	instance service.Instance,
) error {
	return nil
}

func (m *mongoAccountManager) GetUpdater(
	service.Plan,
) (service.Updater, error) {
	// There isn't a need to do any "pre-provision here. just the update step"
	return service.NewUpdater(
		service.NewUpdatingStep("updateARMTemplate", m.updateARMTemplate),
	)
}

func (m *mongoAccountManager) updateARMTemplate(
	_ context.Context,
	instance service.Instance,
) (service.InstanceDetails, error) {
	err := m.cosmosAccountManager.updateDeployment(
		instance.ProvisioningParameters,
		instance.UpdatingParameters,
		instance.Details.(*cosmosdbInstanceDetails),
		"MongoDB",
		"",
		map[string]string{},
	)
	if err != nil {
		return nil, fmt.Errorf("error deploying ARM template: %s", err)
	}
	return instance.Details, nil
}

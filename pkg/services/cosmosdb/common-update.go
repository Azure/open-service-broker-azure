package cosmosdb

import (
	"context"
	"github.com/Azure/open-service-broker-azure/pkg/service"
)

func (
	c *cosmosAccountManager,
) ValidateUpdatingParameters(service.Instance) error {
	return nil
}

func (c *cosmosAccountManager) GetUpdater(service.Plan) (service.Updater, error) {
	// There isn't a need to do any "pre-provision here. just the update step"
	return service.NewUpdater(
		service.NewUpdatingStep("updateARMTemplate", c.updateARMTemplate),
	)
}

func (c *cosmosAccountManager) updateARMTemplate(
	_ context.Context,
	instance service.Instance,
) (service.InstanceDetails, error) {
	return instance.Details, nil
}

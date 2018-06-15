// +build experimental

package cosmosdb

import (
	"github.com/Azure/open-service-broker-azure/pkg/service"
)

func (
	c *cosmosAccountManager,
) ValidateUpdatingParameters(service.Instance) error {
	return nil
}

func (
	c *cosmosAccountManager,
) GetUpdater(service.Plan) (service.Updater, error) {
	return service.NewUpdater()
}

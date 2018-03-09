package cosmosdb

import (
	"github.com/Azure/open-service-broker-azure/pkg/service"
)

func (c *cosmosManager) ValidateUpdatingParameters(service.Instance) error {
	return nil
}

func (c *cosmosManager) GetUpdater(service.Plan) (service.Updater, error) {
	return service.NewUpdater()
}

package cosmosdb

import (
	"github.com/Azure/open-service-broker-azure/pkg/service"
)

func (m *mongoManager) ValidateUpdatingParameters(service.Instance) error {
	return nil
}

func (m *mongoManager) GetUpdater(service.Plan) (service.Updater, error) {
	return service.NewUpdater()
}

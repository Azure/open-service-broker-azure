package cosmosdb

import (
	"github.com/Azure/open-service-broker-azure/pkg/service"
)

func (m *mongoAccountManager) ValidateUpdatingParameters(service.Instance) error {
	return nil
}

func (m *mongoAccountManager) GetUpdater(service.Plan) (service.Updater, error) {
	return service.NewUpdater()
}

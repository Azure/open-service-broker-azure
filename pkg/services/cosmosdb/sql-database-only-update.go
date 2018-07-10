package cosmosdb

import (
	"github.com/Azure/open-service-broker-azure/pkg/service"
)

func (s *sqlDatabaseManager) ValidateUpdatingParameters(
	service.Instance,
) error {
	return nil
}

func (s *sqlDatabaseManager) GetUpdater(service.Plan) (service.Updater, error) {
	return service.NewUpdater()
}

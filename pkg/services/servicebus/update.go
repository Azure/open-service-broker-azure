package servicebus

import (
	"github.com/Azure/open-service-broker-azure/pkg/service"
)

func (s *namespaceManager) ValidateUpdatingParameters(service.Instance) error {
	return nil
}

func (s *namespaceManager) GetUpdater(service.Plan) (service.Updater, error) {
	return service.NewUpdater()
}

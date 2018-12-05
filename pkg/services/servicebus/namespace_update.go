package servicebus

import (
	"github.com/Azure/open-service-broker-azure/pkg/service"
)

func (nm *namespaceManager) ValidateUpdatingParameters(service.Instance) error {
	return nil
}

func (nm *namespaceManager) GetUpdater(service.Plan) (service.Updater, error) {
	return service.NewUpdater()
}

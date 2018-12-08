package servicebus

import (
	"github.com/Azure/open-service-broker-azure/pkg/service"
)

func (qm *queueManager) ValidateUpdatingParameters(service.Instance) error {
	return nil
}

func (qm *queueManager) GetUpdater(service.Plan) (service.Updater, error) {
	return service.NewUpdater()
}

package servicebus

import (
	"github.com/Azure/open-service-broker-azure/pkg/service"
)

func (tm *topicManager) ValidateUpdatingParameters(service.Instance) error {
	return nil
}

func (tm *topicManager) GetUpdater(service.Plan) (service.Updater, error) {
	return service.NewUpdater()
}

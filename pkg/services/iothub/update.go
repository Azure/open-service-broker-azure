package iothub

import (
	"github.com/Azure/open-service-broker-azure/pkg/service"
)

func (i *iotHubManager) ValidateUpdatingParameters(service.Instance) error {
	return nil
}

func (i *iotHubManager) GetUpdater(service.Plan) (service.Updater, error) {
	return service.NewUpdater()
}

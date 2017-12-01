package rediscache

import (
	"github.com/Azure/azure-service-broker/pkg/service"
)

func (s *serviceManager) ValidateUpdatingParameters(
	updatingParameters service.UpdatingParameters,
) error {
	return nil
}

func (s *serviceManager) GetUpdater(service.Plan) (service.Updater, error) {
	return service.NewUpdater()
}

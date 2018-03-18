package acr

import (
	"github.com/Azure/open-service-broker-azure/pkg/service"
)

func (s *serviceManager) ValidateUpdatingParameters(
	updatingParameters service.UpdatingParameters,
) error {
	return nil
}

func (s *serviceManager) GetUpdater(service.Plan) (service.Updater, error) {
	return service.NewUpdater()
}

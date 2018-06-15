// +build experimental

package search

import (
	"open-service-broker-azure/pkg/service"
)

func (s *serviceManager) ValidateUpdatingParameters(service.Instance) error {
	return nil
}

func (s *serviceManager) GetUpdater(service.Plan) (service.Updater, error) {
	return service.NewUpdater()
}

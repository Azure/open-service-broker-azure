package search

import (
	"github.com/Azure/azure-service-broker/pkg/service"
)

func (m *module) ValidateUpdatingParameters(
	updatingParameters service.UpdatingParameters,
) error {
	return nil
}

func (m *module) GetUpdater(string, string) (service.Updater, error) {
	return service.NewUpdater()
}

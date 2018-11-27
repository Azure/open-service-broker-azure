package storage

import (
	"github.com/Azure/open-service-broker-azure/pkg/service"
)

func (s *storageManager) Unbind(service.Instance, service.Binding) error {
	return nil
}

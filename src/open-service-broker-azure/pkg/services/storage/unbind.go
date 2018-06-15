// +build experimental

package storage

import (
	"open-service-broker-azure/pkg/service"
)

func (s *serviceManager) Unbind(service.Instance, service.Binding) error {
	return nil
}

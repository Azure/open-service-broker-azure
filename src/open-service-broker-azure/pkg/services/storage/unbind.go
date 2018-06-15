// +build experimental

package storage

import (
	"github.com/Azure/open-service-broker-azure/pkg/service"
)

func (s *serviceManager) Unbind(service.Instance, service.Binding) error {
	return nil
}

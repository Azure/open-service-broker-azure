// +build experimental

package servicebus

import (
	"open-service-broker-azure/pkg/service"
)

func (s *serviceManager) Unbind(
	_ service.Instance,
	_ service.Binding,
) error {
	return nil
}

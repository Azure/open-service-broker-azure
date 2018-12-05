package servicebus

import (
	"github.com/Azure/open-service-broker-azure/pkg/service"
)

func (tm *topicManager) Unbind(
	_ service.Instance,
	_ service.Binding,
) error {
	return nil
}

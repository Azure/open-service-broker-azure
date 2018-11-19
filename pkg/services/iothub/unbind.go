package iothub

import (
	"github.com/Azure/open-service-broker-azure/pkg/service"
)

func (i *iotHubManager) Unbind(
	_ service.Instance,
	_ service.Binding,
) error {
	return nil
}

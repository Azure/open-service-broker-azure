package storage

import (
	"github.com/Azure/azure-service-broker/pkg/service"
)

func (m *module) Unbind(
	_ service.StandardProvisioningContext,
	_ service.ProvisioningContext,
	_ service.BindingContext,
) error {
	return nil
}

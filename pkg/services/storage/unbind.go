package storage

import (
	"github.com/Azure/azure-service-broker/pkg/service"
)

func (m *module) Unbind(
	provisioningContext service.ProvisioningContext, // nolint: unparam
	bindingContext service.BindingContext,
) error {
	return nil
}

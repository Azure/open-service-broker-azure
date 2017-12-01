package aci

import (
	"github.com/Azure/azure-service-broker/pkg/service"
)

func (s *serviceManager) Unbind(
	_ service.StandardProvisioningContext,
	_ service.ProvisioningContext,
	_ service.BindingContext,
) error {
	return nil
}

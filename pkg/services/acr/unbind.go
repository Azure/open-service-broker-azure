package acr

import "github.com/Azure/open-service-broker-azure/pkg/service"

func (s *serviceManager) Unbind(
	_ service.StandardProvisioningContext,
	_ service.ProvisioningContext,
	_ service.BindingContext,
) error {
	return nil
}

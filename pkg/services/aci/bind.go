package aci

import (
	"fmt"

	"github.com/Azure/open-service-broker-azure/pkg/service"
)

func (s *serviceManager) ValidateBindingParameters(
	bindingParameters service.BindingParameters,
) error {
	// There are no parameters for binding to ACI, so there is nothing
	// to validate
	return nil
}

func (s *serviceManager) Bind(
	instance service.Instance,
	_ service.BindingParameters,
) (service.BindingContext, service.Credentials, error) {
	pc, ok := instance.ProvisioningContext.(*aciProvisioningContext)
	if !ok {
		return nil, nil, fmt.Errorf(
			"error casting instance.ProvisioningContext as *aciProvisioningContext",
		)
	}

	return &aciBindingContext{},
		&aciCredentials{
			PublicIPv4Address: pc.PublicIPv4Address,
		},
		nil
}

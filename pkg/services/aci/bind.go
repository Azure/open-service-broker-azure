package aci

import (
	"fmt"

	"github.com/Azure/azure-service-broker/pkg/service"
)

func (m *module) ValidateBindingParameters(
	bindingParameters service.BindingParameters,
) error {
	// There are no parameters for binding to ACI, so there is nothing
	// to validate
	return nil
}

func (m *module) Bind(
	_ service.StandardProvisioningContext,
	provisioningContext service.ProvisioningContext,
	bindingParameters service.BindingParameters,
) (service.BindingContext, service.Credentials, error) {
	pc, ok := provisioningContext.(*aciProvisioningContext)
	if !ok {
		return nil, nil, fmt.Errorf(
			"error casting provisioningContext as *aciProvisioningContext",
		)
	}

	return &aciBindingContext{},
		&aciCredentials{
			PublicIPv4Address: pc.PublicIPv4Address,
		},
		nil
}

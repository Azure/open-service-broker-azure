package keyvault

import (
	"fmt"

	"github.com/Azure/open-service-broker-azure/pkg/service"
)

func (s *serviceManager) ValidateBindingParameters(
	bindingParameters service.BindingParameters,
) error {
	// There are no parameters for binding to Key vault, so there is nothing
	// to validate
	return nil
}

func (s *serviceManager) Bind(
	instance service.Instance,
	bindingParameters service.BindingParameters,
) (service.BindingContext, service.Credentials, error) {
	pc, ok := instance.ProvisioningContext.(*keyvaultProvisioningContext)
	if !ok {
		return nil, nil, fmt.Errorf(
			"error casting instance.ProvisioningContext as " +
				"*keyvaultProvisioningContext",
		)
	}

	return &keyvaultBindingContext{},
		&Credentials{
			VaultURI:     pc.VaultURI,
			ClientID:     pc.ClientID,
			ClientSecret: pc.ClientSecret,
		},
		nil
}

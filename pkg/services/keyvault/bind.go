package keyvault

import (
	"fmt"

	"github.com/Azure/azure-service-broker/pkg/service"
)

func (m *module) ValidateBindingParameters(
	bindingParameters service.BindingParameters,
) error {
	// There are no parameters for binding to Key vault, so there is nothing
	// to validate
	return nil
}

func (m *module) Bind(
	provisioningContext service.ProvisioningContext,
	bindingParameters service.BindingParameters,
) (service.BindingContext, service.Credentials, error) {
	pc, ok := provisioningContext.(*keyvaultProvisioningContext)
	if !ok {
		return nil, nil, fmt.Errorf(
			"error casting provisioningContext as *keyvaultProvisioningContext",
		)
	}

	return &keyvaultBindingContext{},
		&keyvaultCredentials{
			VaultURI:     pc.VaultURI,
			ClientID:     pc.ClientID,
			ClientSecret: pc.ClientSecret,
		},
		nil
}

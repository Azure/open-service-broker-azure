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
	service.Instance,
	service.BindingParameters,
) (service.BindingDetails, error) {
	return &keyvaultBindingDetails{}, nil
}

func (s *serviceManager) GetCredentials(
	instance service.Instance,
	_ service.Binding,
) (service.Credentials, error) {
	dt, ok := instance.Details.(*keyvaultInstanceDetails)
	if !ok {
		return nil, fmt.Errorf(
			"error casting instance.Details as *keyvaultInstanceDetails",
		)
	}
	return &Credentials{
		VaultURI:     dt.VaultURI,
		ClientID:     dt.ClientID,
		ClientSecret: dt.ClientSecret,
	}, nil
}

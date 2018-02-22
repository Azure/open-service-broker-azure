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
	service.SecureBindingParameters,
) (service.BindingDetails, service.SecureBindingDetails, error) {
	return &keyvaultBindingDetails{}, &keyvaultSecureBindingDetails{}, nil
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
	sdt, ok := instance.SecureDetails.(*keyvaultSecureInstanceDetails)
	if !ok {
		return nil, fmt.Errorf(
			"error casting instance.SecureDetails as *keyvaultSecureInstanceDetails",
		)
	}
	return &Credentials{
		VaultURI:     dt.VaultURI,
		ClientID:     dt.ClientID,
		ClientSecret: sdt.ClientSecret,
	}, nil
}

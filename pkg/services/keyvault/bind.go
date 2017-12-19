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
	_ service.BindingParameters,
) (service.BindingDetails, service.Credentials, error) {
	dt, ok := instance.Details.(*keyvaultInstanceDetails)
	if !ok {
		return nil, nil, fmt.Errorf(
			"error casting instance.Details as *keyvaultInstanceDetails",
		)
	}

	return &keyvaultBindingDetails{},
		&Credentials{
			VaultURI:     dt.VaultURI,
			ClientID:     dt.ClientID,
			ClientSecret: dt.ClientSecret,
		},
		nil
}

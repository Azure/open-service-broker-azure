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
	service.Instance,
	service.BindingParameters,
) (service.BindingDetails, service.SecureBindingDetails, error) {
	return &aciBindingDetails{}, &aciSecureBindingDetails{}, nil
}

func (s *serviceManager) GetCredentials(
	instance service.Instance,
	_ service.Binding,
) (service.Credentials, error) {
	dt, ok := instance.Details.(*aciInstanceDetails)
	if !ok {
		return nil, fmt.Errorf(
			"error casting instance.Details as *aciInstanceDetails",
		)
	}
	return &aciCredentials{
		PublicIPv4Address: dt.PublicIPv4Address,
	}, nil
}

package search

import (
	"fmt"

	"github.com/Azure/open-service-broker-azure/pkg/service"
)

func (s *serviceManager) ValidateBindingParameters(
	bindingParameters service.BindingParameters,
) error {
	// There are no parameters for binding to,
	// so there is nothing to validate
	return nil
}

func (s *serviceManager) Bind(
	service.Instance,
	service.BindingParameters,
) (service.BindingDetails, error) {
	return &searchBindingDetails{}, nil
}

func (s *serviceManager) GetCredentials(
	instance service.Instance,
	_ service.Binding,
) (service.Credentials, error) {
	dt, ok := instance.Details.(*searchInstanceDetails)
	if !ok {
		return nil, fmt.Errorf(
			"error casting instance.Details as *searchInstanceDetails",
		)
	}
	return &searchCredentials{
		ServiceName: dt.ServiceName,
		APIKey:      dt.APIKey,
	}, nil
}

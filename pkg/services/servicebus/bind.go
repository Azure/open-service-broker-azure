package servicebus

import (
	"fmt"

	"github.com/Azure/open-service-broker-azure/pkg/service"
)

func (s *serviceManager) ValidateBindingParameters(
	bindingParameters service.BindingParameters,
) error {
	// There are no parameters for binding to Service Bus, so there is nothing
	// to validate
	return nil
}

func (s *serviceManager) Bind(
	service.Instance,
	service.BindingParameters,
) (service.BindingDetails, error) {
	return &serviceBusBindingDetails{}, nil
}

func (s *serviceManager) GetCredentials(
	instance service.Instance,
	_ service.Binding,
) (service.Credentials, error) {
	pc, ok := instance.Details.(*serviceBusInstanceDetails)
	if !ok {
		return nil, fmt.Errorf(
			"error casting instance.Details as *serviceBusInstanceDetails",
		)
	}
	return &Credentials{
		ConnectionString: pc.ConnectionString,
		PrimaryKey:       pc.PrimaryKey,
	}, nil
}

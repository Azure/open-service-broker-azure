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
	service.SecureBindingParameters,
) (service.BindingDetails, service.SecureBindingDetails, error) {
	return &serviceBusBindingDetails{}, &serviceBusSecureBindingDetails{}, nil
}

func (s *serviceManager) GetCredentials(
	instance service.Instance,
	_ service.Binding,
) (service.Credentials, error) {
	sdt, ok := instance.SecureDetails.(*serviceBusSecureInstanceDetails)
	if !ok {
		return nil, fmt.Errorf(
			"error casting instance.SecureDetails as " +
				"*serviceBusSecureInstanceDetails",
		)
	}
	return &Credentials{
		ConnectionString: sdt.ConnectionString,
		PrimaryKey:       sdt.PrimaryKey,
	}, nil
}

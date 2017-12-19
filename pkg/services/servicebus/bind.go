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
	instance service.Instance,
	_ service.BindingParameters,
) (service.BindingContext, service.Credentials, error) {
	dt, ok := instance.Details.(*serviceBusInstanceDetails)
	if !ok {
		return nil, nil, fmt.Errorf(
			"error casting instance.Details as *serviceBusInstanceDetails",
		)
	}

	return &serviceBusBindingContext{},
		&Credentials{
			ConnectionString: dt.ConnectionString,
			PrimaryKey:       dt.PrimaryKey,
		},
		nil
}

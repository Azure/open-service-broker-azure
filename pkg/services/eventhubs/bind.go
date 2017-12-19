package eventhubs

import (
	"fmt"

	"github.com/Azure/open-service-broker-azure/pkg/service"
)

func (s *serviceManager) ValidateBindingParameters(
	bindingParameters service.BindingParameters,
) error {
	// There are no parameters for binding to Azure Event Hubs, so there is nothing
	// to validate
	return nil
}

func (s *serviceManager) Bind(
	instance service.Instance,
	_ service.BindingParameters,
) (service.BindingDetails, service.Credentials, error) {
	dt, ok := instance.Details.(*eventHubInstanceDetails)
	if !ok {
		return nil, nil, fmt.Errorf(
			"error casting instance.Details as *eventHubInstanceDetails",
		)
	}

	return &eventHubBindingDetails{},
		&Credentials{
			ConnectionString: dt.ConnectionString,
			PrimaryKey:       dt.PrimaryKey,
		},
		nil
}

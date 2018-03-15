package eventhubs

import (
	"github.com/Azure/open-service-broker-azure/pkg/service"
)

func (s *serviceManager) ValidateBindingParameters(
	service.BindingParameters,
	service.SecureBindingParameters,
) error {
	// There are no parameters for binding to Azure Event Hubs, so there is nothing
	// to validate
	return nil
}

func (s *serviceManager) Bind(
	service.Instance,
	service.BindingParameters,
	service.SecureBindingParameters,
) (service.BindingDetails, service.SecureBindingDetails, error) {
	return nil, nil, nil
}

func (s *serviceManager) GetCredentials(
	instance service.Instance,
	_ service.Binding,
) (service.Credentials, error) {
	// In this particular case, we can just return the secure details because
	// its fields do not differ at all from the credentials we wish to return.
	return instance.SecureDetails, nil
}

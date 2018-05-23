package cognitive

import (
	"github.com/Azure/open-service-broker-azure/pkg/service"
)

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
	dt := instanceDetails{}
	if err := service.GetStructFromMap(instance.Details, &dt); err != nil {
		return nil, err
	}
	return &credentials{
		CognitiveKey: dt.CognitiveKey,
		Endpoint:     dt.Endpoint,
	}, nil
}

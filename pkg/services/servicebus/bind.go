package servicebus

import (
	"github.com/Azure/open-service-broker-azure/pkg/service"
)

func (s *namespaceManager) Bind(
	service.Instance,
	service.BindingParameters,
) (service.BindingDetails, error) {
	return nil, nil
}

func (s *namespaceManager) GetCredentials(
	instance service.Instance,
	_ service.Binding,
) (service.Credentials, error) {
	dt := instance.Details.(*instanceDetails)
	return credentials{
		ConnectionString: string(dt.ConnectionString),
		PrimaryKey:       string(dt.PrimaryKey),
	}, nil
}

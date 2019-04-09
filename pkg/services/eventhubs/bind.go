package eventhubs

import (
	"github.com/Azure/open-service-broker-azure/pkg/service"
)

func (s *serviceManager) Bind(
	service.Instance,
	service.BindingParameters,
) (service.BindingDetails, error) {
	return nil, nil
}

func (s *serviceManager) GetCredentials(
	instance service.Instance,
	_ service.Binding,
) (service.Credentials, error) {
	dt := instance.Details.(*instanceDetails)
	return credentials{
		ConnectionString:  string(dt.ConnectionString),
		PrimaryKey:        string(dt.PrimaryKey),
		EventHubNamespace: dt.EventHubNamespace,
		EventHubName:      dt.EventHubName,
	}, nil
}

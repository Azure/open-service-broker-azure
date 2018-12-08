package servicebus

import (
	"github.com/Azure/open-service-broker-azure/pkg/service"
)

func (nm *namespaceManager) Bind(
	service.Instance,
	service.BindingParameters,
) (service.BindingDetails, error) {
	return nil, nil
}

func (nm *namespaceManager) GetCredentials(
	instance service.Instance,
	_ service.Binding,
) (service.Credentials, error) {
	dt := instance.Details.(*namespaceInstanceDetails)
	return namespaceCredentials{
		ConnectionString: string(dt.ConnectionString),
		PrimaryKey:       string(dt.PrimaryKey),
		NamespaceName:    dt.NamespaceName,
	}, nil
}

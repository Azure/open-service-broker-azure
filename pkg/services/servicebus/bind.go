package servicebus

import (
	"fmt"

	"github.com/Azure/azure-service-broker/pkg/service"
)

func (m *module) ValidateBindingParameters(
	bindingParameters service.BindingParameters,
) error {
	// There are no parameters for binding to Service Bus, so there is nothing
	// to validate
	return nil
}

func (m *module) Bind(
	provisioningContext service.ProvisioningContext,
	bindingParameters service.BindingParameters,
) (service.BindingContext, service.Credentials, error) {
	pc, ok := provisioningContext.(*serviceBusProvisioningContext)
	if !ok {
		return nil, nil, fmt.Errorf(
			"error casting provisioningContext as serviceBusProvisioningContext",
		)
	}

	return &serviceBusBindingContext{},
		&Credentials{
			ConnectionString: pc.ConnectionString,
			PrimaryKey:       pc.PrimaryKey,
		},
		nil
}

package postgresqldb

import (
	"fmt"

	"github.com/Azure/open-service-broker-azure/pkg/service"
)

func (d *dbmsOnlyManager) ValidateBindingParameters(
	bindingParameters service.BindingParameters,
) error {
	// There are no parameters for binding to PostgreSQL, so there is nothing
	// to validate
	return nil
}

func (d *dbmsOnlyManager) Bind(
	_ service.Instance,
	_ service.BindingParameters,
) (service.BindingDetails, error) {
	return nil, fmt.Errorf("service is not bindable")
}

func (d *dbmsOnlyManager) GetCredentials(
	_ service.Instance,
	_ service.Binding,
) (service.Credentials, error) {
	return nil, fmt.Errorf("service is not bindable")
}

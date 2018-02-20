package sqldb

import (
	"fmt"

	"github.com/Azure/open-service-broker-azure/pkg/service"
)

func (v *vmOnlyManager) ValidateBindingParameters(
	bindingParameters service.BindingParameters,
) error {
	// There are no parameters for binding to MSSQL, so there is nothing
	// to validate
	return nil
}

//Bind is not valid for VM only,
//TBD behavior
func (v *vmOnlyManager) Bind(
	_ service.Instance,
	_ service.BindingParameters,
) (service.BindingDetails, service.SecureBindingDetails, error) {
	return nil, nil, fmt.Errorf("service is not bindable")
}

func (v *vmOnlyManager) GetCredentials(
	instance service.Instance,
	binding service.Binding,
) (service.Credentials, error) {
	return nil, fmt.Errorf("service is not bindable")
}

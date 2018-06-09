// +build experimental

package mssql

import (
	"fmt"

	"github.com/Azure/open-service-broker-azure/pkg/service"
)

// TODO: Bind is not valid for DBMS only; determine correct behavior
func (d *dbmsManager) Bind(
	service.Instance,
	service.BindingParameters,
	service.SecureBindingParameters,
) (service.BindingDetails, service.SecureBindingDetails, error) {
	return nil, nil, fmt.Errorf("service is not bindable")
}

func (d *dbmsManager) GetCredentials(
	instance service.Instance,
	binding service.Binding,
) (service.Credentials, error) {
	return nil, fmt.Errorf("service is not bindable")
}

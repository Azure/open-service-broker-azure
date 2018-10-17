package mssqldr

import (
	"fmt"

	"github.com/Azure/open-service-broker-azure/pkg/service"
)

// TODO: Bind is not valid for DBMS only; determine correct behavior
func (d *dbmsPairRegisteredManager) Bind(
	service.Instance,
	service.BindingParameters,
) (service.BindingDetails, error) {
	return nil, fmt.Errorf("service is not bindable")
}

func (d *dbmsPairRegisteredManager) GetCredentials(
	instance service.Instance,
	binding service.Binding,
) (service.Credentials, error) {
	return nil, fmt.Errorf("service is not bindable")
}

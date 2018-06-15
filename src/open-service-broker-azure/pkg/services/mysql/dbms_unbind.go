package mysql

import (
	"fmt"

	"open-service-broker-azure/pkg/service"
)

func (v *dbmsManager) Unbind(
	service.Instance,
	service.Binding,
) error {
	return fmt.Errorf("service is not bindable")
}

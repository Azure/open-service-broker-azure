package mysqldb

import (
	"fmt"

	"github.com/Azure/open-service-broker-azure/pkg/service"
)

func (v *dbmsOnlyManager) Unbind(
	_ service.Instance,
	_ service.Binding,
) error {
	return fmt.Errorf("service is not bindable")
}

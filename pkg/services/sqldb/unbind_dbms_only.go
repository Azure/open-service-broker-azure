package sqldb

import (
	"github.com/Azure/open-service-broker-azure/pkg/service"
)

//TODO : Unbind is not valid for VM only.
//Determine what to do.
func (d *dbmsOnlyManager) Unbind(service.Instance, service.Binding) error {
	return nil
}

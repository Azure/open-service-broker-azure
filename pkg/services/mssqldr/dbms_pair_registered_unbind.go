package mssqldr

import (
	"github.com/Azure/open-service-broker-azure/pkg/service"
)

// TODO: Unbind is not valid for DBMS only; determine correct behavior
func (d *dbmsPairRegisteredManager) Unbind(
	service.Instance,
	service.Binding,
) error {
	return nil
}

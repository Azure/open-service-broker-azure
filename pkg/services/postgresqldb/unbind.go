package postgresqldb

import (
	"fmt"

	"github.com/Azure/open-service-broker-azure/pkg/service"
)

func (s *serviceManager) Unbind(
	instance service.Instance,
	bindingContext service.BindingContext,
) error {
	dt, ok := instance.Details.(*postgresqlInstanceDetails)
	if !ok {
		return fmt.Errorf(
			"error casting instance.Details as *postgresqlInstanceDetails",
		)
	}
	bc, ok := bindingContext.(*postgresqlBindingContext)
	if !ok {
		return fmt.Errorf(
			"error casting bindingContext as *postgresqlBindingContext",
		)
	}

	db, err := getDBConnection(dt, primaryDB)
	if err != nil {
		return err
	}
	defer db.Close() // nolint: errcheck

	_, err = db.Exec(
		fmt.Sprintf("drop role %s", bc.LoginName),
	)
	if err != nil {
		return fmt.Errorf(`error dropping role "%s": %s`, bc.LoginName, err)
	}

	return nil
}

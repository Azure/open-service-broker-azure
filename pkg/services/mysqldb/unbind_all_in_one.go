package mysqldb

import (
	"fmt"

	"github.com/Azure/open-service-broker-azure/pkg/service"
)

func (a *allInOneManager) Unbind(
	instance service.Instance,
	bindingDetails service.BindingDetails,
) error {
	dt, ok := instance.Details.(*allInOneMysqlInstanceDetails)
	if !ok {
		return fmt.Errorf(
			"error casting instance.Details as *allInOneMysqlInstanceDetails",
		)
	}
	bc, ok := bindingDetails.(*mysqlBindingDetails)
	if !ok {
		return fmt.Errorf(
			"error casting bindingDetails as *mysqlBindingDetails",
		)
	}

	db, err := a.getDBConnection(dt)
	if err != nil {
		return err
	}
	defer db.Close() // nolint: errcheck

	// Open doesn't open a connection. Validate DSN data:
	err = db.Ping()
	if err != nil {
		return err
	}

	_, err = db.Exec(
		fmt.Sprintf("DROP USER '%s'@'%%'", bc.LoginName),
	)
	if err != nil {
		return fmt.Errorf(
			`error dropping user "%s": %s`,
			bc.LoginName,
			err,
		)
	}

	return nil
}

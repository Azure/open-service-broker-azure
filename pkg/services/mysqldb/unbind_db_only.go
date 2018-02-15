package mysqldb

import (
	"fmt"

	"github.com/Azure/open-service-broker-azure/pkg/service"
)

func (d *dbOnlyManager) Unbind(
	instance service.Instance,
	bindingDetails service.BindingDetails,
) error {
	pdt, ok := instance.Parent.Details.(*dbmsOnlyMysqlInstanceDetails)
	if !ok {
		return fmt.Errorf(
			"error casting instance.Parent.Details as *dbmsOnlyMysqlInstanceDetails",
		)
	}
	spdt, ok := instance.Parent.SecureDetails.(*dbmsOnlyMysqlSecureInstanceDetails)
	if !ok {
		return fmt.Errorf(
			"error casting instance.Parent.SecureDetails as " +
				"*dbmsOnlyMysqlSecureInstanceDetails",
		)
	}
	dt, ok := instance.Details.(*dbOnlyMysqlInstanceDetails)
	if !ok {
		return fmt.Errorf(
			"error casting instance.Details as *dbOnlyMysqlInstanceDetails",
		)
	}
	bc, ok := bindingDetails.(*mysqlBindingDetails)
	if !ok {
		return fmt.Errorf(
			"error casting bindingDetails as *mysqlBindingDetails",
		)
	}

	db, err := d.getDBConnection(pdt, spdt, dt)
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

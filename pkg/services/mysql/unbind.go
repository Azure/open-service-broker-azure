package mysql

import (
	"fmt"

	"github.com/Azure/open-service-broker-azure/pkg/service"
)

func (s *serviceManager) Unbind(
	_ service.StandardProvisioningContext,
	provisioningContext service.ProvisioningContext,
	bindingContext service.BindingContext,
) error {
	pc, ok := provisioningContext.(*mysqlProvisioningContext)
	if !ok {
		return fmt.Errorf(
			"error casting provisioningContext as *mysqlProvisioningContext",
		)
	}
	bc, ok := bindingContext.(*mysqlBindingContext)
	if !ok {
		return fmt.Errorf(
			"error casting bindingContext as *mysqlBindingContext",
		)
	}

	db, err := getDBConnection(pc)
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

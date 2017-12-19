package sqldb

import (
	"fmt"

	"github.com/Azure/open-service-broker-azure/pkg/service"
)

//TODO: What behavior do we want for unbind on a non-bindable service.
func (s *vmServiceManager) Unbind(
	instance service.Instance,
	bindingDetails service.BindingDetails,
) error {
	return nil
}

//TODO implement db only scenario unbind
func (d *dbServiceManager) Unbind(
	_ service.Instance,
	_ service.BindingDetails,
) error {
	return nil
}

func (a *allServiceManager) Unbind(
	instance service.Instance,
	bindingDetails service.BindingDetails,
) error {
	dt, ok := instance.Details.(*mssqlAllInOneInstanceDetails)
	if !ok {
		return fmt.Errorf(
			"error casting instance.Details as  mssqlAllInOneInstanceDetails",
		)
	}
	bc, ok := bindingDetails.(*mssqlBindingDetails)
	if !ok {
		return fmt.Errorf(
			"error casting bindingDetails as *mssqlBindingDetails",
		)
	}

	// connect to new database to drop user for the login
	db, err := getDBConnection(
		dt.AdministratorLogin,
		dt.AdministratorLoginPassword,
		dt.FullyQualifiedDomainName,
		dt.DatabaseName,
	)
	if err != nil {
		return err
	}
	defer db.Close() // nolint: errcheck

	if _, err = db.Exec(
		fmt.Sprintf("DROP USER \"%s\"", bc.LoginName),
	); err != nil {
		return fmt.Errorf(
			`error dropping user "%s": %s`,
			bc.LoginName,
			err,
		)
	}

	// connect to master database to drop login
	masterDb, err := getDBConnection(
		dt.AdministratorLogin,
		dt.AdministratorLoginPassword,
		dt.FullyQualifiedDomainName,
		"master",
	)
	if err != nil {
		return err
	}
	defer masterDb.Close() // nolint: errcheck

	if _, err = masterDb.Exec(
		fmt.Sprintf("DROP LOGIN \"%s\"", bc.LoginName),
	); err != nil {
		return fmt.Errorf(
			`error dropping login "%s": %s`,
			bc.LoginName,
			err,
		)
	}

	return nil
}

package sqldb

import (
	"fmt"

	"github.com/Azure/open-service-broker-azure/pkg/service"
)

func unbind(
	administratorLogin string,
	administratorPassword string,
	fqdn string,
	databaseName string,
	bc *mssqlBindingDetails,
) error {
	// connect to new database to drop user for the login
	db, err := getDBConnection(
		administratorLogin,
		administratorPassword,
		fqdn,
		databaseName)
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
		administratorLogin,
		administratorPassword,
		fqdn,
		"master")
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

func (a *allInOneManager) Unbind(
	instance service.Instance,
	bindingDetails service.BindingDetails,
) error {
	dt, ok := instance.Details.(*mssqlInstanceDetails)
	if !ok {
		return fmt.Errorf(
			"error casting instance.Details as *mssqlInstanceDetails",
		)
	}
	bc, ok := bindingDetails.(*mssqlBindingDetails)
	if !ok {
		return fmt.Errorf(
			"error casting bindingDetails as *mssqlBindingDetails",
		)
	}

	return unbind(
		dt.AdministratorLogin,
		dt.AdministratorLoginPassword,
		dt.FullyQualifiedDomainName,
		dt.DatabaseName,
		bc,
	)
}

//TODO : Unbind is not valid for VM only.
//Determine what to do.
func (v *vmOnlyManager) Unbind(
	_ service.Instance,
	_ service.BindingDetails,
) error {
	return nil
}

func (d *dbOnlyManager) Unbind(
	instance service.Instance,
	bindingDetails service.BindingDetails,
) error {
	dt, ok := instance.Details.(*mssqlInstanceDetails)
	if !ok {
		return fmt.Errorf(
			"error casting instance.Details as *mssqlInstanceDetails",
		)
	}
	bc, ok := bindingDetails.(*mssqlBindingDetails)
	if !ok {
		return fmt.Errorf(
			"error casting bindingDetails as *mssqlBindingDetails",
		)
	}

	return unbind(
		dt.AdministratorLogin,
		dt.AdministratorLoginPassword,
		dt.FullyQualifiedDomainName,
		dt.DatabaseName,
		bc,
	)
}

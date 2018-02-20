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
	binding service.Binding,
) error {
	dt, ok := instance.Details.(*mssqlAllInOneInstanceDetails)
	if !ok {
		return fmt.Errorf(
			"error casting instance.Details as *mssqlAllInOneInstanceDetails",
		)
	}
	sdt, ok := instance.SecureDetails.(*mssqlAllInOneSecureInstanceDetails)
	if !ok {
		return fmt.Errorf(
			"error casting instance.SecureDetails as " +
				"*mssqlAllInOneSecureInstanceDetails",
		)
	}
	bd, ok := binding.Details.(*mssqlBindingDetails)
	if !ok {
		return fmt.Errorf(
			"error casting binding.Details as *mssqlBindingDetails",
		)
	}

	return unbind(
		dt.AdministratorLogin,
		sdt.AdministratorLoginPassword,
		dt.FullyQualifiedDomainName,
		dt.DatabaseName,
		bd,
	)
}

//TODO : Unbind is not valid for VM only.
//Determine what to do.
func (v *vmOnlyManager) Unbind(service.Instance, service.Binding) error {
	return nil
}

func (d *dbOnlyManager) Unbind(
	instance service.Instance,
	binding service.Binding,
) error {
	dt, ok := instance.Details.(*mssqlDBOnlyInstanceDetails)
	if !ok {
		return fmt.Errorf(
			"error casting instance.Details as *mssqlDBOnlyInstanceDetails",
		)
	}
	bd, ok := binding.Details.(*mssqlBindingDetails)
	if !ok {
		return fmt.Errorf(
			"error casting binding.Details as *mssqlBindingDetails",
		)
	}
	//Parent should be set by the framework, but return an error if it is not set.
	if instance.Parent == nil {
		return fmt.Errorf("parent instance not set")
	}
	pdt, ok := instance.Parent.Details.(*mssqlVMOnlyInstanceDetails)
	if !ok {
		return fmt.Errorf(
			"error casting instance.Parent.Details as" +
				"*mssqlVMOnlyInstanceDetails",
		)
	}
	spdt, ok := instance.Parent.SecureDetails.(*mssqlVMOnlySecureInstanceDetails)
	if !ok {
		return fmt.Errorf(
			"error casting instance.Parent.SecureDetails as" +
				"*mssqlVMOnlySecureInstanceDetails",
		)
	}

	return unbind(
		pdt.AdministratorLogin,
		spdt.AdministratorLoginPassword,
		dt.FullyQualifiedDomainName,
		dt.DatabaseName,
		bd,
	)
}

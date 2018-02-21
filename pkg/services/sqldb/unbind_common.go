package sqldb

import (
	"fmt"
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

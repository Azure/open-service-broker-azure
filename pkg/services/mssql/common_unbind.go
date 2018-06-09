// +build experimental

package mssql

import (
	"fmt"
)

func unbind(
	administratorLogin string,
	administratorPassword string,
	fqdn string,
	databaseName string,
	bd bindingDetails,
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
		fmt.Sprintf("DROP USER \"%s\"", bd.LoginName),
	); err != nil {
		return fmt.Errorf(
			`error dropping user "%s": %s`,
			bd.LoginName,
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
		fmt.Sprintf("DROP LOGIN \"%s\"", bd.LoginName),
	); err != nil {
		return fmt.Errorf(
			`error dropping login "%s": %s`,
			bd.LoginName,
			err,
		)
	}

	return nil
}

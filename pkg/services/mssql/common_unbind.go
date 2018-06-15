package mssql

import (
	"fmt"
)

func unbind(
	administratorLogin string,
	administratorPassword string,
	fqdn string,
	databaseName string,
	bd *bindingDetails,
) error {
	// connect to database to drop user
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
		fmt.Sprintf("DROP USER \"%s\"", bd.Username),
	); err != nil {
		return fmt.Errorf(
			`error dropping user "%s": %s`,
			bd.Username,
			err,
		)
	}

	return nil
}

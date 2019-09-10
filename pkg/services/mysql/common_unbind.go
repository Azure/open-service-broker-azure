package mysql

import (
	"fmt"
)

func unbind(
	enforceSSL bool,
	sqlDatabaseDNSSuffix string,
	serverName string,
	administratorLogin string,
	administratorLoginPassword string,
	fullyQualifiedDomainName string,
	databaseName string,
	bindingDetails *bindingDetails,
) error {
	db, err := createDBConnection(
		enforceSSL,
		sqlDatabaseDNSSuffix,
		serverName,
		administratorLogin,
		administratorLoginPassword,
		fullyQualifiedDomainName,
		databaseName,
	)
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
		fmt.Sprintf("DROP USER '%s'@'%%'", bindingDetails.LoginName),
	)
	if err != nil {
		return fmt.Errorf(
			`error dropping user "%s": %s`,
			bindingDetails.LoginName,
			err,
		)
	}
	return nil
}

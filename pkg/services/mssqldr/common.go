package mssqldr

import (
	"database/sql"
	"fmt"
	"net/url"

	_ "github.com/denisenkom/go-mssqldb" // MS SQL Driver
)

func getDBConnection(
	administratorLogin string,
	administratorLoginPassword string,
	fullyQualifiedDomainName string,
	databaseName string,
) (*sql.DB, error) {

	query := url.Values{}
	query.Add("database", databaseName)
	query.Add("encrypt", "true")
	query.Add("TrustServerCertificate", "true")

	u := &url.URL{
		Scheme: "sqlserver",
		User: url.UserPassword(
			administratorLogin,
			administratorLoginPassword,
		),
		Host:     fmt.Sprintf("%s:1433", fullyQualifiedDomainName),
		RawQuery: query.Encode(),
	}

	db, err := sql.Open("mssql", u.String())
	if err != nil {
		return nil, fmt.Errorf("error validating the database arguments: %s", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("error connecting to the database: %s", err)
	}

	return db, nil
}

func validateServerAdmin(
	administratorLogin string,
	administratorLoginPassword string,
	fullyQualifiedDomainName string,
) error {
	// connect to master database
	masterDb, err := getDBConnection(
		administratorLogin,
		administratorLoginPassword,
		fullyQualifiedDomainName,
		"master",
	)
	if err != nil {
		return err
	}
	defer masterDb.Close() // nolint: errcheck

	// Is there a better approach to verify if it is a sys admin?
	rows, err := masterDb.Query("SELECT 1 FROM fn_my_permissions (NULL, 'DATABASE') WHERE permission_name='ALTER ANY USER'") // nolint: lll
	if err != nil {
		return fmt.Errorf(
			"error querying SELECT from table fn_my_permissions: %s",
			err,
		)
	}
	defer rows.Close() // nolint: errcheck
	if !rows.Next() {
		return fmt.Errorf(
			"error user doesn't have permission 'ALTER ANY USER'",
		)
	}
	if err := rows.Err(); err != nil {
		return fmt.Errorf(
			"error iterating rows: %s",
			err,
		)
	}

	return nil
}

package postgresqldb

import (
	"database/sql"
	"fmt"
)

func getDBConnection(
	dt *postgresqlInstanceDetails,
	dbName string,
) (*sql.DB, error) {
	var connectionStrTemplate string
	if dt.EnforceSSL {
		connectionStrTemplate =
			"postgres://postgres@%s:%s@%s/%s?sslmode=require"
	} else {
		connectionStrTemplate = "postgres://postgres@%s:%s@%s/%s"
	}
	db, err := sql.Open("postgres", fmt.Sprintf(
		connectionStrTemplate,
		dt.ServerName,
		dt.AdministratorLoginPassword,
		dt.FullyQualifiedDomainName,
		dbName,
	))
	if err != nil {
		return nil, fmt.Errorf("error connecting to the database: %s", err)
	}
	return db, err
}

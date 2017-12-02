package postgresqldb

import (
	"database/sql"
	"fmt"
)

func getDBConnection(
	pc *postgresqlProvisioningContext,
	dbName string,
) (*sql.DB, error) {
	var connectionStrTemplate string
	if pc.EnforceSSL {
		connectionStrTemplate =
			"postgres://postgres@%s:%s@%s/%s?sslmode=require"
	} else {
		connectionStrTemplate = "postgres://postgres@%s:%s@%s/%s"
	}
	db, err := sql.Open("postgres", fmt.Sprintf(
		connectionStrTemplate,
		pc.ServerName,
		pc.AdministratorLoginPassword,
		pc.FullyQualifiedDomainName,
		dbName,
	))
	if err != nil {
		return nil, fmt.Errorf("error connecting to the database: %s", err)
	}
	return db, err
}

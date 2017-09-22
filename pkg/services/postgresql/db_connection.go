package postgresql

import (
	"database/sql"
	"fmt"
)

func getDBConnection(pc *postgresqlProvisioningContext) (*sql.DB, error) {
	var connectionStrTemplate string
	if pc.EnforceSSL {
		connectionStrTemplate =
			"postgres://postgres@%s:%s@%s/postgres?sslmode=require"
	} else {
		connectionStrTemplate = "postgres://postgres@%s:%s@%s/postgres"
	}
	db, err := sql.Open("postgres", fmt.Sprintf(
		connectionStrTemplate,
		pc.ServerName,
		pc.AdministratorLoginPassword,
		pc.FullyQualifiedDomainName,
	))
	if err != nil {
		return nil, fmt.Errorf("error connecting to the database: %s", err)
	}
	return db, err
}

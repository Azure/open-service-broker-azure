package postgresql

import (
	"database/sql"
	"fmt"
)

func getDBConnection(pc *postgresqlProvisioningContext) (*sql.DB, error) {
	db, err := sql.Open("postgres", fmt.Sprintf(
		"postgres://postgres@%s:%s@%s/postgres?sslmode=require",
		pc.ServerName,
		pc.AdministratorLoginPassword,
		pc.FullyQualifiedDomainName,
	))
	if err != nil {
		return nil, fmt.Errorf("error connecting to the database: %s", err)
	}
	return db, err
}

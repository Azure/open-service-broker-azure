package postgresql

import (
	"database/sql"
	"fmt"
)

const primaryDB = "postgres"

func getDBConnection(
	enforceSSL bool,
	serverName string,
	administratorLoginPassword string,
	fullyQualifiedDomainName string,
	dbName string,
) (*sql.DB, error) {
	var connectionStrTemplate string
	if enforceSSL {
		connectionStrTemplate =
			"postgres://postgres@%s:%s@%s/%s?sslmode=require"
	} else {
		connectionStrTemplate = "postgres://postgres@%s:%s@%s/%s"
	}
	db, err := sql.Open("postgres", fmt.Sprintf(
		connectionStrTemplate,
		serverName,
		administratorLoginPassword,
		fullyQualifiedDomainName,
		dbName,
	))
	if err != nil {
		return nil, fmt.Errorf("error connecting to the database: %s", err)
	}
	return db, err
}

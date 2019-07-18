package postgresql

import (
	"database/sql"
	"fmt"

	"github.com/Azure/open-service-broker-azure/pkg/service"
	_ "github.com/lib/pq" // Postgres SQL driver
)

const primaryDB = "postgres"

var dbExtensionsSchema = &service.ArrayPropertySchema{
	Title:       "Database extensions",
	Description: "Database extensions to install",
	ItemsSchema: &service.StringPropertySchema{
		Title:       "Name",
		Description: "Extension Name",
	},
}

func getDBConnection(
	enforceSSL bool,
	administratorLogin string,
	serverName string,
	administratorLoginPassword string,
	fullyQualifiedDomainName string,
	dbName string,
) (*sql.DB, error) {
	var connectionStrTemplate string
	if enforceSSL {
		connectionStrTemplate =
			"postgres://%s@%s:%s@%s/%s?sslmode=require"
	} else {
		connectionStrTemplate = "postgres://%s@%s:%s@%s/%s"
	}
	db, err := sql.Open("postgres", fmt.Sprintf(
		connectionStrTemplate,
		administratorLogin,
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

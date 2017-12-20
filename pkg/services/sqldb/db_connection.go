package sqldb

import (
	"database/sql"
	"fmt"
	"net/url"

	_ "github.com/denisenkom/go-mssqldb" // MS SQL Driver
)

func getDBConnection(
	administratorLogin string,
	administratorPassword string,
	fqdn string,
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
			administratorPassword,
		),
		Host:     fmt.Sprintf("%s:1433", fqdn),
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

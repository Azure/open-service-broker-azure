package mssql

import (
	"fmt"
	"net/url"

	"github.com/Azure/open-service-broker-azure/pkg/generate"
	"github.com/Azure/open-service-broker-azure/pkg/service"
	log "github.com/Sirupsen/logrus"
)

func bind(
	administratorLogin string,
	administratorPassword string,
	fqdn string,
	databaseName string,
) (service.BindingDetails, error) {

	username := generate.NewIdentifier()
	password := generate.NewPassword()

	// connect to new database to create user
	db, err := getDBConnection(
		administratorLogin,
		administratorPassword,
		fqdn,
		databaseName,
	)
	if err != nil {
		return nil, err
	}
	defer db.Close() // nolint: errcheck

	tx, err := db.Begin()
	if err != nil {
		return nil, fmt.Errorf(
			"error starting transaction on the new database: %s",
			err,
		)
	}
	defer func() {
		if err != nil {
			if err = tx.Rollback(); err != nil {
				log.WithField("error", err).
					Error("error rolling back transaction on the new database")
			}
		}
	}()
	if _, err = tx.Exec(
		fmt.Sprintf("CREATE USER \"%s\" WITH PASSWORD='%s'", username, password),
	); err != nil {
		return nil, fmt.Errorf(
			`error creating user "%s": %s`,
			username,
			err,
		)
	}
	if _, err = tx.Exec(
		fmt.Sprintf("GRANT CONTROL to \"%s\"", username),
	); err != nil {
		return nil, fmt.Errorf(
			`error granting CONTROL to user "%s": %s`,
			username,
			err,
		)
	}
	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf(
			"error committing transaction on the new database: %s",
			err,
		)
	}
	return &bindingDetails{
		Username: username,
		Password: service.SecureString(password),
	}, nil
}

func createCredential(
	fqdn string,
	database string,
	username string,
	password string,
) credentials {

	port := 1433

	jdbcTemplate := "jdbc:sqlserver://%s:%d;database=%s;user=%s;" +
		"password=%s;encrypt=true;trustServerCertificate=true;"

	jdbc := fmt.Sprintf(
		jdbcTemplate,
		fqdn,
		port,
		database,
		username,
		password,
	)

	uriTemplate :=
		"sqlserver://%s:%s@%s:%d/%s;encrypt=true;trustServerCertificate=true"

	uri := fmt.Sprintf(
		uriTemplate,
		url.QueryEscape(username),
		password,
		fqdn,
		port,
		database,
	)
	return credentials{
		Host:     fqdn,
		Port:     port,
		Database: database,
		Username: username,
		Password: password,
		JDBC:     jdbc,
		URI:      uri,
		Encrypt:  true,
	}
}

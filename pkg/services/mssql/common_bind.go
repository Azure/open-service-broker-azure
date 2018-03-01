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
) (service.BindingDetails, service.SecureBindingDetails, error) {

	loginName := generate.NewIdentifier()
	password := generate.NewPassword()

	// connect to master database to create login
	masterDb, err := getDBConnection(
		administratorLogin,
		administratorPassword,
		fqdn,
		"master",
	)
	if err != nil {
		return nil, nil, err
	}
	defer masterDb.Close() // nolint: errcheck

	if _, err = masterDb.Exec(
		fmt.Sprintf("CREATE LOGIN \"%s\" WITH PASSWORD='%s'", loginName, password),
	); err != nil {
		return nil, nil, fmt.Errorf(
			`error creating login "%s": %s`,
			loginName,
			err,
		)
	}

	// connect to new database to create user for the login
	db, err := getDBConnection(
		administratorLogin,
		administratorPassword,
		fqdn,
		databaseName,
	)
	if err != nil {
		return nil, nil, err
	}
	defer db.Close() // nolint: errcheck

	tx, err := db.Begin()
	if err != nil {
		return nil, nil, fmt.Errorf(
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
			// Drop the login created in the last step
			if _, err = masterDb.Exec(
				fmt.Sprintf("DROP LOGIN \"%s\"", loginName),
			); err != nil {
				log.WithField("error", err).
					Error("error dropping login on master database")
			}
		}
	}()
	if _, err = tx.Exec(
		fmt.Sprintf("CREATE USER \"%s\" FOR LOGIN \"%s\"", loginName, loginName),
	); err != nil {
		return nil, nil, fmt.Errorf(
			`error creating user "%s": %s`,
			loginName,
			err,
		)
	}
	if _, err = tx.Exec(
		fmt.Sprintf("GRANT CONTROL to \"%s\"", loginName),
	); err != nil {
		return nil, nil, fmt.Errorf(
			`error granting CONTROL to user "%s": %s`,
			loginName,
			err,
		)
	}
	if err = tx.Commit(); err != nil {
		return nil, nil, fmt.Errorf(
			"error committing transaction on the new database: %s",
			err,
		)
	}

	return &bindingDetails{
			LoginName: loginName,
		},
		&secureBindingDetails{
			Password: password,
		},
		nil
}

func createCredential(
	fqdn string,
	database string,
	username string,
	password string,
) *Credentials {

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
	return &Credentials{
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
